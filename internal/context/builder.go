package context

import (
	"fmt"
	"strings"

	"aether/model"
	"github.com/pkoukk/tiktoken-go"
)

// PreparedContext represents structured context ready for LLM consumption.
type PreparedContext struct {
	SystemPrompt string       `json:"system_prompt"`
	UserPrompt   string       `json:"user_prompt"`
	TokenCount   int          `json:"token_count"`
	Entities     []model.Node `json:"entities"`
}

// Builder formats entity nodes and graph relationships into token-budgeted prompt text.
type Builder struct {
	tokenBpe *tiktoken.Tiktoken
}

// NewBuilder creates a Context Builder with a tiktoken encoding for token estimation.
func NewBuilder(modelName string) (*Builder, error) {
	if modelName == "" {
		modelName = "gpt-4o"
	}
	bpe, err := tiktoken.EncodingForModel(modelName)
	if err != nil {
		bpe, _ = tiktoken.GetEncoding("cl100k_base")
	}
	return &Builder{tokenBpe: bpe}, nil
}

// EstimateTokens calculates the approximate token count for a text string.
func (b *Builder) EstimateTokens(text string) int {
	if b.tokenBpe == nil {
		return len(strings.Fields(text)) // Fallback rough word count
	}
	tokens := b.tokenBpe.Encode(text, nil, nil)
	return len(tokens)
}

// BuildStructuredContext formats a list of graph nodes and relationships into a token-limited context payload.
func (b *Builder) BuildStructuredContext(systemRole string, goal string, nodes []model.Node, tokenBudget int) PreparedContext {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("## Project Goal\n%s\n\n", goal))
	sb.WriteString("## Project Context (Knowledge Model Entities)\n")

	usedTokens := b.EstimateTokens(systemRole + sb.String())
	var includedNodes []model.Node

	for _, node := range nodes {
		nodeText := fmt.Sprintf("- [%s] %s (%s): line %d-%d\n", node.Type, node.Name, node.Path, node.LineStart, node.LineEnd)
		nodeTokens := b.EstimateTokens(nodeText)

		if usedTokens+nodeTokens > tokenBudget {
			sb.WriteString(fmt.Sprintf("\n[Truncated context: %d remaining entities omitted to respect token limit]\n", len(nodes)-len(includedNodes)))
			break
		}

		sb.WriteString(nodeText)
		usedTokens += nodeTokens
		includedNodes = append(includedNodes, node)
	}

	return PreparedContext{
		SystemPrompt: systemRole,
		UserPrompt:   sb.String(),
		TokenCount:   usedTokens,
		Entities:     includedNodes,
	}
}
