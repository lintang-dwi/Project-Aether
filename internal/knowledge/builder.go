package knowledge

import (
	"context"
	"fmt"

	"aether/internal/parser"
	"aether/model"
)

// Builder builds KnowledgeModel from file parse results.
type Builder struct {
	pipeline *parser.Pipeline
}

// NewBuilder creates a new Knowledge Model Builder.
func NewBuilder(pipeline *parser.Pipeline) *Builder {
	return &Builder{pipeline: pipeline}
}

// BuildFromResults aggregates multiple ParseResult structs into a KnowledgeModel.
func (b *Builder) BuildFromResults(results []*parser.ParseResult) *KnowledgeModel {
	km := NewKnowledgeModel()

	for _, res := range results {
		for _, entity := range res.Entities {
			km.AddNode(entity)
		}
		for _, rel := range res.Relationships {
			km.AddEdge(rel)
		}

		// Connect file node to imported dependencies
		fileNodeID := fmt.Sprintf("file:%s", res.FilePath)
		for _, imp := range res.Imports {
			importNodeID := fmt.Sprintf("import:%s", imp)
			km.AddNode(model.Node{
				ID:       importNodeID,
				Type:     model.EntityModule,
				Name:     imp,
				Language: res.Language,
			})
			km.AddEdge(model.Edge{
				ID:     fmt.Sprintf("edge:imports:%s:%s", fileNodeID, importNodeID),
				FromID: fileNodeID,
				ToID:   importNodeID,
				Type:   model.RelImports,
			})
		}
	}

	return km
}

// ProcessFile parses content and merges it into an existing KnowledgeModel.
func (b *Builder) ProcessFile(ctx context.Context, km *KnowledgeModel, filePath string, content []byte) error {
	res, err := b.pipeline.ParseFile(ctx, filePath, content)
	if err != nil {
		return err
	}
	merged := b.BuildFromResults([]*parser.ParseResult{res})
	for id, node := range merged.Entities {
		km.Entities[id] = node
	}
	for id, edge := range merged.Relationships {
		km.Relationships[id] = edge
	}
	return nil
}
