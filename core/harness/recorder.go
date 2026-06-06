package harness

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/wzhongyou/llmgate/core"
)

type ctxKeySkipRecord struct{}

// SkipRecordCtx returns a context that tells Recorder to skip recording.
func SkipRecordCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKeySkipRecord{}, true)
}

// Record is one recorded request-response pair.
type Record struct {
	Timestamp time.Time          `json:"timestamp"`
	Provider  string             `json:"provider,omitempty"`
	Model     string             `json:"model,omitempty"`
	LatencyMs float64            `json:"latency_ms,omitempty"`
	Request   *core.ChatRequest  `json:"request"`
	Response  *core.ChatResponse `json:"response,omitempty"`
	Error     string             `json:"error,omitempty"`
}

// Recorder implements core.Hook, writing request/response pairs to a JSONL file.
type Recorder struct {
	mu      sync.Mutex
	file    *os.File
	enc     *json.Encoder
	count   int64
	enabled bool
}

func NewRecorder(path string) (*Recorder, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &Recorder{file: f, enc: json.NewEncoder(f), enabled: true}, nil
}

func (r *Recorder) AfterChat(ctx context.Context, req *core.ChatRequest, resp *core.ChatResponse, err error) {
	if ctx.Value(ctxKeySkipRecord{}) != nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.enabled {
		return
	}
	rec := Record{
		Timestamp: time.Now(),
		Request:   req,
		Response:  resp,
	}
	if resp != nil {
		rec.Provider = resp.Provider
		rec.Model = resp.Model
		rec.LatencyMs = float64(resp.Latency.Microseconds()) / 1000.0
	}
	if err != nil {
		rec.Error = err.Error()
	}
	r.enc.Encode(rec)
	r.count++
}

func (r *Recorder) SetEnabled(on bool) {
	r.mu.Lock()
	r.enabled = on
	r.mu.Unlock()
}

func (r *Recorder) Enabled() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.enabled
}

func (r *Recorder) Count() int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.count
}

func (r *Recorder) Path() string {
	return r.file.Name()
}

func (r *Recorder) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.file.Close()
}
