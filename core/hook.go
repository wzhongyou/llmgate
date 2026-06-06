package core

import "context"

// Hook observes engine-level events without modifying control flow.
type Hook interface {
	AfterChat(ctx context.Context, req *ChatRequest, resp *ChatResponse, err error)
}

// AddHook registers a hook that will be called after each Chat invocation.
func (e *Engine) AddHook(h Hook) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.hooks = append(e.hooks, h)
}

func (e *Engine) fireAfterChat(ctx context.Context, req *ChatRequest, resp *ChatResponse, err error) {
	e.mu.RLock()
	hooks := e.hooks
	e.mu.RUnlock()
	for _, h := range hooks {
		h.AfterChat(ctx, req, resp, err)
	}
}
