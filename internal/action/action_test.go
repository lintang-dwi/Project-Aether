package action_test

import (
	"os"
	"path/filepath"
	"testing"

	"aether/internal/action"
	"aether/internal/eventbus"
)

func TestActionProcessor_CreateAndEdit(t *testing.T) {
	tmpDir := t.TempDir()
	eb := eventbus.NewEventBus(100)
	defer eb.Close()

	proc := action.NewProcessor(tmpDir, eb)

	// Create
	createOp := &action.CreateFileOp{
		Path:    "hello.txt",
		Content: "Hello World",
	}
	resCreate := proc.ExecuteOperation(createOp)
	if !resCreate.Success {
		t.Fatalf("create failed: %s", resCreate.Error)
	}

	content, _ := os.ReadFile(filepath.Join(tmpDir, "hello.txt"))
	if string(content) != "Hello World" {
		t.Errorf("expected 'Hello World', got '%s'", string(content))
	}

	// Edit
	editOp := &action.EditFileOp{
		Path:       "hello.txt",
		NewContent: "Hello GraphOS",
	}
	resEdit := proc.ExecuteOperation(editOp)
	if !resEdit.Success {
		t.Fatalf("edit failed: %s", resEdit.Error)
	}

	contentUpdated, _ := os.ReadFile(filepath.Join(tmpDir, "hello.txt"))
	if string(contentUpdated) != "Hello GraphOS" {
		t.Errorf("expected 'Hello GraphOS', got '%s'", string(contentUpdated))
	}
}
