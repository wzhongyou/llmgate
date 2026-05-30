package core

import "encoding/json"

type oaiMsg struct {
	Role             string      `json:"role"`
	Content          interface{} `json:"content"`
	ToolCalls        []ToolCall  `json:"tool_calls,omitempty"`
	ToolCallID       string      `json:"tool_call_id,omitempty"`
	ReasoningContent string      `json:"reasoning_content,omitempty"`
}

type oaiContentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

func OpenAIMessages(req *ChatRequest) []oaiMsg {
	var msgs []oaiMsg
	if req.System != "" {
		msgs = append(msgs, oaiMsg{Role: "system", Content: req.System})
	}
	for _, m := range req.Messages {
		om := oaiMsg{Role: m.Role, ToolCallID: m.ToolCallID, ReasoningContent: m.ReasoningContent}
		if len(m.ToolCalls) > 0 {
			om.ToolCalls = m.ToolCalls
			if m.Content != "" {
				om.Content = m.Content
			}
		} else if len(m.ContentParts) > 0 {
			parts := make([]oaiContentPart, len(m.ContentParts))
			for i, p := range m.ContentParts {
				parts[i] = oaiContentPart{Type: p.Type, Text: p.Text, ImageURL: p.ImageURL}
			}
			om.Content = parts
		} else {
			om.Content = m.Content
		}
		msgs = append(msgs, om)
	}
	return msgs
}

// OpenAIBody builds the request body map for OpenAI-compatible /chat/completions.
func OpenAIBody(model string, stream bool, req *ChatRequest) map[string]interface{} {
	body := map[string]interface{}{
		"model":    model,
		"messages": OpenAIMessages(req),
	}
	if stream {
		body["stream"] = true
	}
	if req.MaxTokens != nil {
		body["max_tokens"] = *req.MaxTokens
	}
	if req.Temperature != nil {
		body["temperature"] = *req.Temperature
	}
	if len(req.Tools) > 0 {
		body["tools"] = req.Tools
		if req.ToolChoice != nil {
			body["tool_choice"] = req.ToolChoice
		}
	}
	if req.TopP != nil {
		body["top_p"] = *req.TopP
	}
	if len(req.Stop) > 0 {
		body["stop"] = req.Stop
	}
	if req.FrequencyPenalty != nil {
		body["frequency_penalty"] = *req.FrequencyPenalty
	}
	if req.PresencePenalty != nil {
		body["presence_penalty"] = *req.PresencePenalty
	}
	if req.Seed != nil {
		body["seed"] = *req.Seed
	}
	if req.ResponseFormat != nil {
		body["response_format"] = req.ResponseFormat
	}
	if req.ThinkingType == "disabled" {
		body["thinking"] = map[string]string{"type": "disabled"}
	}
	return body
}

// OpenAIParseChat parses an OpenAI-compatible chat completion response body.
func OpenAIParseChat(data []byte, providerName string) (*ChatResponse, error) {
	var r struct {
		Choices []struct {
			Message struct {
				Content          *string    `json:"content"`
				ToolCalls        []ToolCall `json:"tool_calls"`
				ReasoningContent string     `json:"reasoning_content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Model string `json:"model"`
		Usage struct {
			PromptTokens      int `json:"prompt_tokens"`
			CompletionTokens  int `json:"completion_tokens"`
			TotalTokens       int `json:"total_tokens"`
			CompletionDetails *struct {
				ReasoningTokens int `json:"reasoning_tokens"`
			} `json:"completion_tokens_details"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, &ProviderError{Provider: providerName, Message: err.Error(), Cause: err}
	}
	if len(r.Choices) == 0 {
		return nil, &ProviderError{Provider: providerName, Message: "no choices in response"}
	}
	ch := r.Choices[0]
	var content string
	if ch.Message.Content != nil {
		content = *ch.Message.Content
	}
	if content == "" && len(ch.Message.ToolCalls) == 0 {
		return nil, &ProviderError{Provider: providerName, Message: "empty response"}
	}
	reasoning := 0
	if r.Usage.CompletionDetails != nil {
		reasoning = r.Usage.CompletionDetails.ReasoningTokens
	}
	return &ChatResponse{
		Content:          content,
		ToolCalls:        ch.Message.ToolCalls,
		FinishReason:     ch.FinishReason,
		Model:            r.Model,
		ReasoningContent: ch.Message.ReasoningContent,
		Usage: Usage{
			InputTokens:     r.Usage.PromptTokens,
			OutputTokens:    r.Usage.CompletionTokens,
			ReasoningTokens: reasoning,
			TotalTokens:     r.Usage.TotalTokens,
		},
	}, nil
}
