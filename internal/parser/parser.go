package parser

import (
	"context"

	"aether/model"
)

// SymbolKind represents the type of code construct extracted from AST.
type SymbolKind string

const (
	KindPackage   SymbolKind = "PACKAGE"
	KindImport    SymbolKind = "IMPORT"
	KindFunction  SymbolKind = "FUNCTION"
	KindMethod    SymbolKind = "METHOD"
	KindStruct    SymbolKind = "STRUCT"
	KindClass     SymbolKind = "CLASS"
	KindInterface SymbolKind = "INTERFACE"
	KindVariable  SymbolKind = "VARIABLE"
)

// Symbol represents an extracted code symbol with location and signature.
type Symbol struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Kind      SymbolKind        `json:"kind"`
	Path      string            `json:"path"`
	LineStart int               `json:"line_start"`
	LineEnd   int               `json:"line_end"`
	Signature string            `json:"signature,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// ParseResult holds the parsed symbols and dependencies extracted from a file.
type ParseResult struct {
	FilePath     string         `json:"file_path"`
	Language     string         `json:"language"`
	Symbols      []Symbol       `json:"symbols"`
	Imports      []string       `json:"imports"`
	Calls        []string       `json:"calls"`
	Entities     []model.Node   `json:"entities"`
	Relationships []model.Edge `json:"relationships"`
}

// LanguageParser defines the interface for language-specific Tree-sitter parsers.
type LanguageParser interface {
	Language() string
	CanParse(filePath string) bool
	Parse(ctx context.Context, filePath string, content []byte) (*ParseResult, error)
}

// Pipeline manages multi-language parsing.
type Pipeline struct {
	parsers map[string]LanguageParser
}

// NewPipeline creates a new multi-language Parser Pipeline.
func NewPipeline() *Pipeline {
	return &Pipeline{
		parsers: make(map[string]LanguageParser),
	}
}

// RegisterParser registers a language parser into the pipeline.
func (p *Pipeline) RegisterParser(parser LanguageParser) {
	p.parsers[parser.Language()] = parser
}

// ParseFile parses a file using the appropriate language parser.
func (p *Pipeline) ParseFile(ctx context.Context, filePath string, content []byte) (*ParseResult, error) {
	for _, parser := range p.parsers {
		if parser.CanParse(filePath) {
			return parser.Parse(ctx, filePath, content)
		}
	}
	// Fallback empty result for unsupported file types
	return &ParseResult{
		FilePath: filePath,
		Language: "unknown",
	}, nil
}
