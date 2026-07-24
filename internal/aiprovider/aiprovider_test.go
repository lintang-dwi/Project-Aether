package aiprovider_test

import (
	"context"
	"testing"

	"aether/internal/aiprovider"
)

type MockProvider struct {
	Response string
}

func (m *MockProvider) Name() string {
	return "mock"
}

func (m *MockProvider) Complete(ctx context.Context, req aiprovider.CompletionRequest) (*aiprovider.CompletionResponse, error) {
	return &aiprovider.CompletionResponse{
		Content:      m.Response,
		FinishReason: "stop",
		Model:        "mock-model",
		TokensUsed:   50,
	}, nil
}

func (m *MockProvider) StreamComplete(ctx context.Context, req aiprovider.CompletionRequest, handler func(chunk string) error) error {
	return handler(m.Response)
}

func TestMockAIProvider(t *testing.T) {
	mock := &MockProvider{Response: "Hello from Mock AI"}

	resp, err := mock.Complete(context.Background(), aiprovider.CompletionRequest{
		Prompt: "Say hi",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Content != "Hello from Mock AI" {
		t.Errorf("expected 'Hello from Mock AI', got '%s'", resp.Content)
	}
	if resp.TokensUsed != 50 {
		t.Errorf("expected 50 tokens used, got %d", resp.TokensUsed)
	}
}
