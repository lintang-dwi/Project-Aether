package context_test

import (
	"fmt"
	"testing"

	"aether/internal/context"
	"aether/internal/graph"
	"aether/internal/knowledge"
	"aether/model"
)

func TestContextEngine_TokenBudgeting(t *testing.T) {
	ge := graph.NewEngine()
	km := knowledge.NewKnowledgeModel()

	// Add 50 nodes
	for i := 1; i <= 50; i++ {
		km.AddNode(model.Node{
			ID:        fmt.Sprintf("node-%d", i),
			Type:      model.EntityFunction,
			Name:      fmt.Sprintf("FunctionName%d", i),
			Path:      fmt.Sprintf("internal/pkg/file%d.go", i),
			LineStart: 10,
			LineEnd:   25,
		})
	}
	_ = ge.LoadKnowledgeModel(km)

	engine, err := context.NewEngine(ge, "gpt-4o")
	if err != nil {
		t.Fatalf("failed to create context engine: %v", err)
	}

	// Budget of 200 tokens
	prepared := engine.PrepareContextForTarget("Refactor internal package", "internal/", 200)

	if prepared.TokenCount > 250 {
		t.Errorf("expected token count close to budget 200, got %d", prepared.TokenCount)
	}
	if len(prepared.Entities) >= 50 {
		t.Errorf("expected entity list to be truncated by token budget, got %d", len(prepared.Entities))
	}
}
