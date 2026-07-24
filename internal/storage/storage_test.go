package storage_test

import (
	"context"
	"testing"

	"aether/internal/storage"
	"aether/model"
)

func TestStorage_SQLiteAndGraphCRUD(t *testing.T) {
	db, err := storage.OpenDB(":memory:")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// 1. Test Metadata Repo
	metaRepo := storage.NewMetadataRepository(db)
	if err := metaRepo.Set(ctx, "workspace_name", "Project Aether"); err != nil {
		t.Fatalf("failed to set metadata: %v", err)
	}

	val, err := metaRepo.Get(ctx, "workspace_name")
	if err != nil {
		t.Fatalf("failed to get metadata: %v", err)
	}
	if val != "Project Aether" {
		t.Errorf("expected 'Project Aether', got '%s'", val)
	}

	// 2. Test Graph Repo
	graphRepo := storage.NewGraphRepository(db)
	node1 := model.Node{
		ID:       "node-1",
		Type:     model.EntityFile,
		Name:     "main.go",
		Path:     "cmd/aether/main.go",
		Language: "go",
	}
	node2 := model.Node{
		ID:       "node-2",
		Type:     model.EntityPackage,
		Name:     "config",
		Path:     "internal/config",
		Language: "go",
	}

	if err := graphRepo.SaveNode(ctx, node1); err != nil {
		t.Fatalf("failed to save node1: %v", err)
	}
	if err := graphRepo.SaveNode(ctx, node2); err != nil {
		t.Fatalf("failed to save node2: %v", err)
	}

	n1Fetched, err := graphRepo.GetNode(ctx, "node-1")
	if err != nil {
		t.Fatalf("failed to fetch node1: %v", err)
	}
	if n1Fetched.Name != "main.go" || n1Fetched.Type != model.EntityFile {
		t.Errorf("unexpected node data: %+v", n1Fetched)
	}

	// Edge
	edge := model.Edge{
		ID:     "edge-1",
		FromID: "node-1",
		ToID:   "node-2",
		Type:   model.RelImports,
		Weight: 1.0,
	}

	if err := graphRepo.SaveEdge(ctx, edge); err != nil {
		t.Fatalf("failed to save edge: %v", err)
	}

	edges, err := graphRepo.GetEdgesFrom(ctx, "node-1")
	if err != nil {
		t.Fatalf("failed to get edges from node-1: %v", err)
	}
	if len(edges) != 1 || edges[0].ToID != "node-2" {
		t.Errorf("unexpected edges result: %+v", edges)
	}
}
