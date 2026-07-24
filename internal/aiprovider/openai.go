package aiprovider

import (
	"context"
	"fmt"
	"io"

	openai "github.com/sashabaranov/go-openai"
)

// OpenAIProvider implements Provider using the official/compatible OpenAI Go SDK.
type OpenAIProvider struct {
	client *openai.Client
	model  string
}

// NewOpenAIProvider initializes an OpenAI provider with API key and default model.
func NewOpenAIProvider(apiKey, defaultModel string) *OpenAIProvider {
	if defaultModel == "" {
		defaultModel = openai.GPT4oMini
	}
	client := openai.NewClient(apiKey)
	return &OpenAIProvider{
		client: client,
		model:  defaultModel,
	}
}

func (p *OpenAIProvider) Name() string {
	return "openai"
}

func (p *OpenAIProvider) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	modelName := req.Model
	if modelName == "" {
		modelName = p.model
	}

	messages := make([]openai.ChatCompletionMessage, 0)
	if req.System != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: req.System,
		})
	}
	if len(req.Messages) > 0 {
		for _, m := range req.Messages {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    m.Role,
				Content: m.Content,
			})
		}
	} else if req.Prompt != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: req.Prompt,
		})
	}

	resp, err := p.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       modelName,
		Messages:    messages,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
	})
	if err != nil {
		return nil, fmt.Errorf("openai chat completion failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("openai returned empty choices")
	}

	return &CompletionResponse{
		Content:      resp.Choices[0].Message.Content,
		FinishReason: string(resp.Choices[0].FinishReason),
		TokensUsed:   resp.Usage.TotalTokens,
		Model:        resp.Model,
	}, nil
}

func (p *OpenAIProvider) StreamComplete(ctx context.Context, req CompletionRequest, handler func(chunk string) error) error {
	modelName := req.Model
	if modelName == "" {
		modelName = p.model
	}

	messages := make([]openai.ChatCompletionMessage, 0)
	if req.System != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: req.System,
		})
	}
	if req.Prompt != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: req.Prompt,
		})
	}

	stream, err := p.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:       modelName,
		Messages:    messages,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
	})
	if err != nil {
		return fmt.Errorf("failed to create openai stream: %w", err)
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errorsIsEOF(err) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error reading stream chunk: %w", err)
		}
		if len(response.Choices) > 0 {
			if err := handler(response.Choices[0].Delta.Content); err != nil {
				return err
			}
		}
	}
}

func errorsIsEOF(err error) bool {
	return err != nil && err.Error() == io.EOF.Error()
}
