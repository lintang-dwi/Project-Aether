package git_test

import (
	"os"
	"path/filepath"
	"testing"

	gitengine "aether/internal/git"
)

func TestGitEngine_CheckpointAndRollback(t *testing.T) {
	tmpDir := t.TempDir()

	engine, err := gitengine.NewEngine(tmpDir)
	if err != nil {
		t.Fatalf("failed to init git engine: %v", err)
	}

	// 1. Initial file & Checkpoint 1
	file1 := filepath.Join(tmpDir, "file1.txt")
	_ = os.WriteFile(file1, []byte("Version 1"), 0644)

	cp1, err := engine.CreateCheckpoint("Initial commit")
	if err != nil {
		t.Fatalf("failed to create checkpoint 1: %v", err)
	}

	// 2. Modify file & Checkpoint 2
	_ = os.WriteFile(file1, []byte("Version 2 - Corrupted"), 0644)

	// 3. Perform Rollback to Checkpoint 1
	if err := engine.Rollback(cp1.Hash); err != nil {
		t.Fatalf("failed to rollback: %v", err)
	}

	contentAfter, _ := os.ReadFile(file1)
	if string(contentAfter) != "Version 1" {
		t.Errorf("expected content 'Version 1' after rollback, got '%s'", string(contentAfter))
	}
}
