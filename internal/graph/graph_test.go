package graph_test

import (
	"testing"

	"aether/internal/graph"
	"aether/internal/knowledge"
	"aether/model"
)

func TestGraphEngine_TraversalAndQuery(t *testing.T) {
	km := knowledge.NewKnowledgeModel()

	// Add nodes: FileA -> FileB -> FileC
	nodeA := model.Node{ID: "file:A.go", Type: model.EntityFile, Name: "A.go", Path: "cmd/A.go", Language: "go"}
	nodeB := model.Node{ID: "file:B.go", Type: model.EntityFile, Name: "B.go", Path: "internal/B.go", Language: "go"}
	nodeC := model.Node{ID: "file:C.go", Type: model.EntityFile, Name: "C.go", Path: "internal/C.go", Language: "go"}

	km.AddNode(nodeA)
	km.AddNode(nodeB)
	km.AddNode(nodeC)

	// Edges: A imports B, B imports C
	km.AddEdge(model.Edge{ID: "e1", FromID: "file:A.go", ToID: "file:B.go", Type: model.RelImports})
	km.AddEdge(model.Edge{ID: "e2", FromID: "file:B.go", ToID: "file:C.go", Type: model.RelImports})

	engine := graph.NewEngine()
	if err := engine.LoadKnowledgeModel(km); err != nil {
		t.Fatalf("failed to load knowledge model: %v", err)
	}

	if engine.VertexCount() != 3 {
		t.Errorf("expected 3 vertices, got %d", engine.VertexCount())
	}

	// Direct dependencies of A.go -> should be B.go
	deps, err := engine.GetDependencies("file:A.go")
	if err != nil {
		t.Fatalf("failed to get dependencies of A.go: %v", err)
	}
	if len(deps) != 1 || deps[0].ID != "file:B.go" {
		t.Errorf("expected dependency [file:B.go], got %v", deps)
	}

	// Dependents of C.go -> should be B.go
	dependents, err := engine.GetDependents("file:C.go")
	if err != nil {
		t.Fatalf("failed to get dependents of C.go: %v", err)
	}
	if len(dependents) != 1 || dependents[0].ID != "file:B.go" {
		t.Errorf("expected dependent [file:B.go], got %v", dependents)
	}

	// Shortest Path A.go -> C.go
	path, err := engine.FindShortestPath("file:A.go", "file:C.go")
	if err != nil {
		t.Fatalf("failed to find path A -> C: %v", err)
	}
	if len(path) != 3 || path[0] != "file:A.go" || path[2] != "file:C.go" {
		t.Errorf("unexpected shortest path: %v", path)
	}

	// QueryNodes
	matches := engine.QueryNodes(graph.FilterOptions{
		PathPrefix: "internal/",
	})
	if len(matches) != 2 {
		t.Errorf("expected 2 nodes under internal/, got %d", len(matches))
	}
}
