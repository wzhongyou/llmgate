package harness

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/wzhongyou/llmgate/core"
)

// ShadowRecord is the result of a shadow request.
type ShadowRecord struct {
	Timestamp time.Time          `json:"timestamp"`
	Provider  string             `json:"provider"`
	Model     string             `json:"model,omitempty"`
	LatencyMs float64            `json:"latency_ms"`
	Request   *core.ChatRequest  `json:"request"`
	Response  *core.ChatResponse `json:"response,omitempty"`
	Error     string             `json:"error,omitempty"`
}

// Shadow implements core.Hook, asynchronously sending each request to a shadow provider.
type Shadow struct {
	mu       sync.Mutex
	engine   *core.Engine
	provider string
	model    string
	file     *os.File
	enc      *json.Encoder
	enabled  bool
}

func NewShadow(engine *core.Engine, provider, model, path string) (*Shadow, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &Shadow{
		engine:   engine,
		provider: provider,
		model:    model,
		file:     f,
		enc:      json.NewEncoder(f),
		enabled:  true,
	}, nil
}

func (s *Shadow) AfterChat(_ context.Context, req *core.ChatRequest, _ *core.ChatResponse, origErr error) {
	s.mu.Lock()
	enabled := s.enabled
	s.mu.Unlock()
	if !enabled || origErr != nil || req == nil {
		return
	}
	go s.dispatch(req)
}

func (s *Shadow) dispatch(req *core.ChatRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	shadowReq := *req
	if s.model != "" {
		shadowReq.Model = s.model
	}

	start := time.Now()
	resp, err := s.engine.ChatWithProvider(ctx, &shadowReq, s.provider)
	latency := float64(time.Since(start).Microseconds()) / 1000.0

	rec := ShadowRecord{
		Timestamp: time.Now(),
		Provider:  s.provider,
		LatencyMs: latency,
		Request:   req,
		Response:  resp,
	}
	if resp != nil {
		rec.Model = resp.Model
	}
	if err != nil {
		rec.Error = err.Error()
	}

	s.mu.Lock()
	s.enc.Encode(rec)
	s.mu.Unlock()
}

func (s *Shadow) SetEnabled(on bool) {
	s.mu.Lock()
	s.enabled = on
	s.mu.Unlock()
}

func (s *Shadow) Enabled() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.enabled
}

func (s *Shadow) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.file.Close()
}
