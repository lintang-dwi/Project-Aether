package parser

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"aether/model"
	sitter "github.com/smacker/go-tree-sitter"
	sitterjs "github.com/smacker/go-tree-sitter/javascript"
)

// TSParser implements LanguageParser for TypeScript & JavaScript using Tree-sitter.
type TSParser struct{}

func NewTSParser() *TSParser {
	return &TSParser{}
}

func (p *TSParser) Language() string {
	return "typescript"
}

func (p *TSParser) CanParse(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".ts" || ext == ".tsx" || ext == ".js" || ext == ".jsx"
}

func (p *TSParser) Parse(ctx context.Context, filePath string, content []byte) (*ParseResult, error) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitterjs.GetLanguage()
	parser.SetLanguage(lang)

	tree, err := parser.ParseCtx(ctx, nil, content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JS/TS AST for %s: %w", filePath, err)
	}
	defer tree.Close()

	result := &ParseResult{
		FilePath: filePath,
		Language: "typescript",
	}

	fileNodeID := fmt.Sprintf("file:%s", filePath)
	result.Entities = append(result.Entities, model.Node{
		ID:       fileNodeID,
		Type:     model.EntityFile,
		Name:     filepath.Base(filePath),
		Path:     filePath,
		Language: "typescript",
	})

	root := tree.RootNode()
	p.extractSymbols(root, content, filePath, fileNodeID, result)

	return result, nil
}

func (p *TSParser) extractSymbols(node *sitter.Node, content []byte, filePath, fileNodeID string, result *ParseResult) {
	if node == nil {
		return
	}

	count := node.ChildCount()
	for i := 0; i < int(count); i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}

		switch child.Type() {
		case "import_statement":
			p.parseImport(child, content, filePath, result)
		case "function_declaration":
			p.parseFunction(child, content, filePath, fileNodeID, result)
		case "class_declaration":
			p.parseClass(child, content, filePath, fileNodeID, result)
		}
	}
}

func (p *TSParser) parseImport(node *sitter.Node, content []byte, filePath string, result *ParseResult) {
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil && child.Type() == "string" {
			impPath := strings.Trim(child.Content(content), `"'`)
			result.Imports = append(result.Imports, impPath)

			result.Symbols = append(result.Symbols, Symbol{
				ID:        fmt.Sprintf("import:%s:%s", filePath, impPath),
				Name:      impPath,
				Kind:      KindImport,
				Path:      filePath,
				LineStart: int(child.StartPoint().Row) + 1,
				LineEnd:   int(child.EndPoint().Row) + 1,
			})
		}
	}
}

func (p *TSParser) parseFunction(node *sitter.Node, content []byte, filePath, fileNodeID string, result *ParseResult) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}
	funcName := nameNode.Content(content)
	funcID := fmt.Sprintf("fn:%s:%s", filePath, funcName)

	sym := Symbol{
		ID:        funcID,
		Name:      funcName,
		Kind:      KindFunction,
		Path:      filePath,
		LineStart: int(node.StartPoint().Row) + 1,
		LineEnd:   int(node.EndPoint().Row) + 1,
	}
	result.Symbols = append(result.Symbols, sym)

	result.Entities = append(result.Entities, model.Node{
		ID:        funcID,
		Type:      model.EntityFunction,
		Name:      funcName,
		Path:      filePath,
		Language:  "typescript",
		LineStart: sym.LineStart,
		LineEnd:   sym.LineEnd,
	})

	result.Relationships = append(result.Relationships, model.Edge{
		ID:     fmt.Sprintf("edge:contains:%s:%s", fileNodeID, funcID),
		FromID: fileNodeID,
		ToID:   funcID,
		Type:   model.RelContains,
	})
}

func (p *TSParser) parseClass(node *sitter.Node, content []byte, filePath, fileNodeID string, result *ParseResult) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}
	className := nameNode.Content(content)
	classID := fmt.Sprintf("class:%s:%s", filePath, className)

	sym := Symbol{
		ID:        classID,
		Name:      className,
		Kind:      KindClass,
		Path:      filePath,
		LineStart: int(node.StartPoint().Row) + 1,
		LineEnd:   int(node.EndPoint().Row) + 1,
	}
	result.Symbols = append(result.Symbols, sym)

	result.Entities = append(result.Entities, model.Node{
		ID:        classID,
		Type:      model.EntityClass,
		Name:      className,
		Path:      filePath,
		Language:  "typescript",
		LineStart: sym.LineStart,
		LineEnd:   sym.LineEnd,
	})

	result.Relationships = append(result.Relationships, model.Edge{
		ID:     fmt.Sprintf("edge:contains:%s:%s", fileNodeID, classID),
		FromID: fileNodeID,
		ToID:   classID,
		Type:   model.RelContains,
	})
}
