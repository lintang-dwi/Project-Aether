package workspace_test

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"aether/internal/eventbus"
	"aether/internal/workspace"
)

func TestWorkspaceScanner_GitignoreFilter(t *testing.T) {
	tmpDir := t.TempDir()

	// Create files
	_ = os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main"), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("# Test"), 0644)
	_ = os.MkdirAll(filepath.Join(tmpDir, "node_modules", "pkg"), 0755)
	_ = os.WriteFile(filepath.Join(tmpDir, "node_modules", "pkg", "index.js"), []byte(""), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, "ignored.log"), []byte("log"), 0644)

	// Create .gitignore
	_ = os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte("*.log\n"), 0644)

	scanner, err := workspace.NewScanner(tmpDir, nil)
	if err != nil {
		t.Fatalf("failed to create scanner: %v", err)
	}

	files, err := scanner.Scan()
	if err != nil {
		t.Fatalf("failed to scan workspace: %v", err)
	}

	fileMap := make(map[string]bool)
	for _, f := range files {
		fileMap[f.Path] = true
	}

	if !fileMap["main.go"] {
		t.Error("expected 'main.go' to be discovered")
	}
	if !fileMap["README.md"] {
		t.Error("expected 'README.md' to be discovered")
	}
	if fileMap["ignored.log"] {
		t.Error("expected 'ignored.log' to be ignored by .gitignore")
	}
	if fileMap["node_modules/pkg/index.js"] {
		t.Error("expected node_modules to be ignored")
	}
}

func TestWorkspaceEngine_WatcherEvents(t *testing.T) {
	tmpDir := t.TempDir()
	eb := eventbus.NewEventBus(100)
	defer eb.Close()

	engine, err := workspace.NewEngine(tmpDir, nil, eb)
	if err != nil {
		t.Fatalf("failed to create workspace engine: %v", err)
	}

	ctx := context.Background()
	if err := engine.Start(ctx); err != nil {
		t.Fatalf("failed to start engine: %v", err)
	}
	defer engine.Stop(ctx)

	var mu sync.Mutex
	var receivedEvent bool
	eb.Subscribe("workspace.*", func(e eventbus.Event) {
		mu.Lock()
		receivedEvent = true
		mu.Unlock()
	})

	// Create a new file in workspace
	testFile := filepath.Join(tmpDir, "new_file.txt")
	_ = os.WriteFile(testFile, []byte("hello"), 0644)

	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	got := receivedEvent
	mu.Unlock()

	if !got {
		t.Error("expected file creation event to be captured by workspace watcher")
	}
}
