package workspace

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"aether/internal/eventbus"
	"github.com/fsnotify/fsnotify"
)

// Watcher monitors workspace file events and publishes them to the EventBus.
type Watcher struct {
	mu         sync.RWMutex
	rootPath   string
	watcher    *fsnotify.Watcher
	scanner    *Scanner
	eventBus   *eventbus.EventBus
	isWatching bool
	cancel     context.CancelFunc
}

// NewWatcher initializes a file watcher for the workspace.
func NewWatcher(rootPath string, scanner *Scanner, eb *eventbus.EventBus) (*Watcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		rootPath: rootPath,
		watcher:  fw,
		scanner:  scanner,
		eventBus: eb,
	}, nil
}

// Start begins listening for file change events.
func (w *Watcher) Start(ctx context.Context) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.isWatching {
		return nil
	}

	// Add root recursively
	err := filepath.Walk(w.rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		relPath, _ := filepath.Rel(w.rootPath, path)
		if relPath != "." && w.scanner != nil && w.scanner.IsIgnored(relPath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.IsDir() {
			return w.watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	watchCtx, cancel := context.WithCancel(ctx)
	w.cancel = cancel
	w.isWatching = true

	go w.loop(watchCtx)
	return nil
}

func (w *Watcher) loop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			_ = err
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			relPath, err := filepath.Rel(w.rootPath, event.Name)
			if err != nil || (w.scanner != nil && w.scanner.IsIgnored(relPath)) {
				continue
			}

			topic := "workspace.file_changed"
			if event.Op&fsnotify.Create == fsnotify.Create {
				topic = "workspace.file_created"
				if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
					_ = w.watcher.Add(event.Name)
				}
			} else if event.Op&fsnotify.Remove == fsnotify.Remove {
				topic = "workspace.file_deleted"
			}

			w.eventBus.Publish(eventbus.NewEvent(topic, "workspace_watcher", map[string]interface{}{
				"path": filepath.ToSlash(relPath),
				"op":   event.Op.String(),
			}))
		}
	}
}

// Stop halts the file watcher.
func (w *Watcher) Stop() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.isWatching {
		return nil
	}

	if w.cancel != nil {
		w.cancel()
	}
	w.isWatching = false
	return w.watcher.Close()
}
