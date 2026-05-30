package core

import "context"

type EmbeddingProvider interface {
	Embed(ctx context.Context, req *EmbedRequest) (*EmbedResponse, error)
}

type EmbedRequest struct {
	Input []string `json:"input"`
	Model string   `json:"model,omitempty"`
}

type EmbedResponse struct {
	Embeddings [][]float32 `json:"embeddings"`
	Model      string      `json:"model"`
	Usage      Usage       `json:"usage"`
}
