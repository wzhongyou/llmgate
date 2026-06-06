package harness

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/wzhongyou/llmgate/core"
)

// Violation represents a single format compliance failure.
type Violation struct {
	Timestamp time.Time `json:"timestamp"`
	Provider  string    `json:"provider"`
	Model     string    `json:"model"`
	Type      string    `json:"type"` // "invalid_json" | "schema_mismatch" | "unknown_finish_reason"
	Detail    string    `json:"detail"`
}

// Probe implements core.Hook, checking response format compliance.
type Probe struct {
	mu         sync.Mutex
	violations []Violation
	maxSize    int
}

func NewProbe(maxSize int) *Probe {
	if maxSize <= 0 {
		maxSize = 100
	}
	return &Probe{maxSize: maxSize}
}

func (p *Probe) AfterChat(_ context.Context, req *core.ChatRequest, resp *core.ChatResponse, err error) {
	if err != nil || resp == nil {
		return
	}

	// Check tool_calls arguments are valid JSON
	for _, tc := range resp.ToolCalls {
		if tc.Function.Arguments == "" {
			continue
		}
		if !json.Valid([]byte(tc.Function.Arguments)) {
			p.record(Violation{
				Timestamp: time.Now(),
				Provider:  resp.Provider,
				Model:     resp.Model,
				Type:      "invalid_json",
				Detail:    "tool_call " + tc.Function.Name + " arguments is not valid JSON",
			})
		}
	}

	// Check finish_reason is in known set
	switch resp.FinishReason {
	case "", "stop", "tool_calls", "length", "content_filter":
		// valid
	default:
		p.record(Violation{
			Timestamp: time.Now(),
			Provider:  resp.Provider,
			Model:     resp.Model,
			Type:      "unknown_finish_reason",
			Detail:    resp.FinishReason,
		})
	}
}

func (p *Probe) record(v Violation) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.violations) >= p.maxSize {
		p.violations = p.violations[1:]
	}
	p.violations = append(p.violations, v)
}

func (p *Probe) Violations() []Violation {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]Violation, len(p.violations))
	copy(out, p.violations)
	return out
}
