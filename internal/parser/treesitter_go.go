package parser

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"aether/model"
	sitter "github.com/smacker/go-tree-sitter"
	sittergolang "github.com/smacker/go-tree-sitter/golang"
)

// GoParser implements LanguageParser for the Go language using Tree-sitter.
type GoParser struct{}

func NewGoParser() *GoParser {
	return &GoParser{}
}

func (p *GoParser) Language() string {
	return "go"
}

func (p *GoParser) CanParse(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".go"
}

func (p *GoParser) Parse(ctx context.Context, filePath string, content []byte) (*ParseResult, error) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sittergolang.GetLanguage()
	parser.SetLanguage(lang)

	tree, err := parser.ParseCtx(ctx, nil, content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Go AST for %s: %w", filePath, err)
	}
	defer tree.Close()

	result := &ParseResult{
		FilePath: filePath,
		Language: "go",
	}

	// Root File Node
	fileNodeID := fmt.Sprintf("file:%s", filePath)
	result.Entities = append(result.Entities, model.Node{
		ID:       fileNodeID,
		Type:     model.EntityFile,
		Name:     filepath.Base(filePath),
		Path:     filePath,
		Language: "go",
	})

	root := tree.RootNode()
	p.extractSymbols(root, content, filePath, fileNodeID, result)

	return result, nil
}

func (p *GoParser) extractSymbols(node *sitter.Node, content []byte, filePath, fileNodeID string, result *ParseResult) {
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
		case "package_clause":
			p.parsePackage(child, content, filePath, fileNodeID, result)
		case "import_declaration":
			p.parseImports(child, content, filePath, fileNodeID, result)
		case "function_declaration":
			p.parseFunction(child, content, filePath, fileNodeID, result)
		case "method_declaration":
			p.parseMethod(child, content, filePath, fileNodeID, result)
		case "type_declaration":
			p.parseType(child, content, filePath, fileNodeID, result)
		}
	}
}

func (p *GoParser) parsePackage(node *sitter.Node, content []byte, filePath, fileNodeID string, result *ParseResult) {
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil && child.Type() == "package_identifier" {
			pkgName := child.Content(content)
			sym := Symbol{
				ID:        fmt.Sprintf("pkg:%s", pkgName),
				Name:      pkgName,
				Kind:      KindPackage,
				Path:      filePath,
				LineStart: int(node.StartPoint().Row) + 1,
				LineEnd:   int(node.EndPoint().Row) + 1,
			}
			result.Symbols = append(result.Symbols, sym)
			break
		}
	}
}

func (p *GoParser) parseImports(node *sitter.Node, content []byte, filePath, fileNodeID string, result *ParseResult) {
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}
		if child.Type() == "import_spec" || child.Type() == "import_spec_list" {
			for j := 0; j < int(child.ChildCount()); j++ {
				spec := child.Child(j)
				if spec != nil && spec.Type() == "interpreted_string_literal" {
					impPath := strings.Trim(spec.Content(content), `"`)
					result.Imports = append(result.Imports, impPath)

					sym := Symbol{
						ID:        fmt.Sprintf("import:%s:%s", filePath, impPath),
						Name:      impPath,
						Kind:      KindImport,
						Path:      filePath,
						LineStart: int(spec.StartPoint().Row) + 1,
						LineEnd:   int(spec.EndPoint().Row) + 1,
					}
					result.Symbols = append(result.Symbols, sym)
				}
			}
		}
	}
}

func (p *GoParser) parseFunction(node *sitter.Node, content []byte, filePath, fileNodeID string, result *ParseResult) {
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
		Language:  "go",
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

func (p *GoParser) parseMethod(node *sitter.Node, content []byte, filePath, fileNodeID string, result *ParseResult) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}
	methodName := nameNode.Content(content)
	methodID := fmt.Sprintf("method:%s:%s", filePath, methodName)

	sym := Symbol{
		ID:        methodID,
		Name:      methodName,
		Kind:      KindMethod,
		Path:      filePath,
		LineStart: int(node.StartPoint().Row) + 1,
		LineEnd:   int(node.EndPoint().Row) + 1,
	}
	result.Symbols = append(result.Symbols, sym)

	result.Entities = append(result.Entities, model.Node{
		ID:        methodID,
		Type:      model.EntityFunction,
		Name:      methodName,
		Path:      filePath,
		Language:  "go",
		LineStart: sym.LineStart,
		LineEnd:   sym.LineEnd,
	})

	result.Relationships = append(result.Relationships, model.Edge{
		ID:     fmt.Sprintf("edge:contains:%s:%s", fileNodeID, methodID),
		FromID: fileNodeID,
		ToID:   methodID,
		Type:   model.RelContains,
	})
}

func (p *GoParser) parseType(node *sitter.Node, content []byte, filePath, fileNodeID string, result *ParseResult) {
	specNode := node.ChildByFieldName("spec")
	if specNode == nil {
		for i := 0; i < int(node.ChildCount()); i++ {
			c := node.Child(i)
			if c != nil && (c.Type() == "type_spec" || c.Type() == "type_alias") {
				specNode = c
				break
			}
		}
	}
	if specNode == nil {
		return
	}

	nameNode := specNode.ChildByFieldName("name")
	if nameNode == nil {
		return
	}
	typeName := nameNode.Content(content)
	typeID := fmt.Sprintf("type:%s:%s", filePath, typeName)

	kind := KindStruct
	typeDefNode := specNode.ChildByFieldName("type")
	if typeDefNode != nil && typeDefNode.Type() == "interface_type" {
		kind = KindInterface
	}

	sym := Symbol{
		ID:        typeID,
		Name:      typeName,
		Kind:      kind,
		Path:      filePath,
		LineStart: int(node.StartPoint().Row) + 1,
		LineEnd:   int(node.EndPoint().Row) + 1,
	}
	result.Symbols = append(result.Symbols, sym)

	entityType := model.EntityStruct
	if kind == KindInterface {
		entityType = model.EntityInterface
	}

	result.Entities = append(result.Entities, model.Node{
		ID:        typeID,
		Type:      entityType,
		Name:      typeName,
		Path:      filePath,
		Language:  "go",
		LineStart: sym.LineStart,
		LineEnd:   sym.LineEnd,
	})

	result.Relationships = append(result.Relationships, model.Edge{
		ID:     fmt.Sprintf("edge:contains:%s:%s", fileNodeID, typeID),
		FromID: fileNodeID,
		ToID:   typeID,
		Type:   model.RelContains,
	})
}
