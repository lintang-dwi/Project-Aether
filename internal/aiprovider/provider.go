package aiprovider

import (
	"context"
)

// CompletionRequest represents a prompt payload for AI models.
type CompletionRequest struct {
	Model       string    `json:"model"`
	Prompt      string    `json:"prompt"`
	System      string    `json:"system,omitempty"`
	Temperature float32   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Messages    []Message `json:"messages,omitempty"`
}

// Message represents a chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CompletionResponse represents the AI model's response.
type CompletionResponse struct {
	Content      string `json:"content"`
	FinishReason string `json:"finish_reason"`
	TokensUsed   int    `json:"tokens_used"`
	Model        string `json:"model"`
}

// Provider defines the unified interface for AI model integrations (Cloud & Local).
type Provider interface {
	Name() string
	Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)
	StreamComplete(ctx context.Context, req CompletionRequest, handler func(chunk string) error) error
}
