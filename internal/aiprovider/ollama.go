package aiprovider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OllamaProvider implements Provider for local models via the Ollama REST API.
type OllamaProvider struct {
	baseURL    string
	defaultModel string
	client     *http.Client
}

type ollamaChatRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
}

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaChatResponse struct {
	Model      string        `json:"model"`
	Message    ollamaMessage `json:"message"`
	Done       bool          `json:"done"`
	TotalDuration int64      `json:"total_duration,omitempty"`
}

// NewOllamaProvider initializes an Ollama provider.
func NewOllamaProvider(baseURL, defaultModel string) *OllamaProvider {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	if defaultModel == "" {
		defaultModel = "llama3.2"
	}
	return &OllamaProvider{
		baseURL:      baseURL,
		defaultModel: defaultModel,
		client:       &http.Client{Timeout: 120 * time.Second},
	}
}

func (p *OllamaProvider) Name() string {
	return "ollama"
}

func (p *OllamaProvider) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	modelName := req.Model
	if modelName == "" {
		modelName = p.defaultModel
	}

	messages := make([]ollamaMessage, 0)
	if req.System != "" {
		messages = append(messages, ollamaMessage{Role: "system", Content: req.System})
	}
	if len(req.Messages) > 0 {
		for _, m := range req.Messages {
			messages = append(messages, ollamaMessage{Role: m.Role, Content: m.Content})
		}
	} else if req.Prompt != "" {
		messages = append(messages, ollamaMessage{Role: "user", Content: req.Prompt})
	}

	body := ollamaChatRequest{
		Model:    modelName,
		Messages: messages,
		Stream:   false,
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ollama request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/api/chat", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ollama API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var chatResp ollamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("failed to decode ollama response: %w", err)
	}

	return &CompletionResponse{
		Content:      chatResp.Message.Content,
		FinishReason: "stop",
		Model:        chatResp.Model,
	}, nil
}

func (p *OllamaProvider) StreamComplete(ctx context.Context, req CompletionRequest, handler func(chunk string) error) error {
	// Fallback to non-stream complete for simplicity
	res, err := p.Complete(ctx, req)
	if err != nil {
		return err
	}
	return handler(res.Content)
}
