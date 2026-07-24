package knowledge_test

import (
	"context"
	"testing"

	"aether/internal/knowledge"
	"aether/internal/parser"
)

func TestKnowledgeBuilder(t *testing.T) {
	pipeline := parser.NewPipeline()
	pipeline.RegisterParser(parser.NewGoParser())

	builder := knowledge.NewBuilder(pipeline)
	km := knowledge.NewKnowledgeModel()

	goCode := []byte(`
package main
import "os"
func Execute() {}
`)

	err := builder.ProcessFile(context.Background(), km, "cmd/app/main.go", goCode)
	if err != nil {
		t.Fatalf("failed to process file: %v", err)
	}

	nodes := km.GetNodes()
	if len(nodes) < 3 { // File, Function Execute, Import os
		t.Errorf("expected at least 3 nodes in KnowledgeModel, got %d", len(nodes))
	}

	edges := km.GetEdges()
	if len(edges) < 2 { // File contains Function, File imports os
		t.Errorf("expected at least 2 edges in KnowledgeModel, got %d", len(edges))
	}
}
