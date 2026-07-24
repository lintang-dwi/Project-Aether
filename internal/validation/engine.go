package validation

import (
	"context"
	"fmt"

	"aether/internal/parser"
)

// ValidationResult represents validation status and diagnostics.
type ValidationResult struct {
	IsValid     bool     `json:"is_valid"`
	Diagnostics []string `json:"diagnostics,omitempty"`
}

// Engine validates code modifications before commit.
type Engine struct {
	pipeline *parser.Pipeline
}

// NewEngine creates a new Validation Engine instance.
func NewEngine(pipeline *parser.Pipeline) *Engine {
	return &Engine{pipeline: pipeline}
}

// ValidateFile checks if the content of a file parses cleanly.
func (e *Engine) ValidateFile(ctx context.Context, filePath string, content []byte) ValidationResult {
	if e.pipeline == nil {
		return ValidationResult{IsValid: true}
	}

	res, err := e.pipeline.ParseFile(ctx, filePath, content)
	if err != nil {
		return ValidationResult{
			IsValid:     false,
			Diagnostics: []string{fmt.Sprintf("Parse error: %v", err)},
		}
	}

	if res.Language != "unknown" && len(res.Entities) == 0 && len(content) > 50 {
		return ValidationResult{
			IsValid:     false,
			Diagnostics: []string{"Empty parse tree produced for non-empty file"},
		}
	}

	return ValidationResult{
		IsValid: true,
	}
}
