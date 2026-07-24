package workspace

import (
	"context"
	"fmt"
	"sync"

	"aether/internal/coordinator"
	"aether/internal/eventbus"
	"aether/model"
)

// Engine implements coordinator.Service for managing the project workspace.
type Engine struct {
	mu       sync.RWMutex
	rootPath string
	scanner  *Scanner
	watcher  *Watcher
	eventBus *eventbus.EventBus
	status   coordinator.ServiceStatus
}

// NewEngine creates a new Workspace Engine service.
func NewEngine(rootPath string, customIgnores []string, eb *eventbus.EventBus) (*Engine, error) {
	scanner, err := NewScanner(rootPath, customIgnores)
	if err != nil {
		return nil, fmt.Errorf("failed to create scanner: %w", err)
	}

	watcher, err := NewWatcher(scanner.GetRootPath(), scanner, eb)
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	return &Engine{
		rootPath: scanner.GetRootPath(),
		scanner:  scanner,
		watcher:  watcher,
		eventBus: eb,
		status:   coordinator.StatusStopped,
	}, nil
}

func (e *Engine) Name() string {
	return "workspace_engine"
}

func (e *Engine) Start(ctx context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.watcher.Start(ctx); err != nil {
		e.status = coordinator.StatusFailed
		return fmt.Errorf("failed to start workspace watcher: %w", err)
	}

	e.status = coordinator.StatusRunning
	return nil
}

func (e *Engine) Stop(ctx context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.watcher.Stop(); err != nil {
		e.status = coordinator.StatusFailed
		return fmt.Errorf("failed to stop workspace watcher: %w", err)
	}

	e.status = coordinator.StatusStopped
	return nil
}

func (e *Engine) Health() coordinator.ServiceStatus {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.status
}

// ScanWorkspace scans the workspace and returns file metadata list.
func (e *Engine) ScanWorkspace() ([]model.FileMetadata, error) {
	return e.scanner.Scan()
}

// GetScanner returns the internal Scanner.
func (e *Engine) GetScanner() *Scanner {
	return e.scanner
}
