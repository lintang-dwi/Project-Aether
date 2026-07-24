package plugin

import (
	"context"
	"fmt"
	"sync"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// Manager manages WebAssembly plugin modules using wazero (pure Go, CGO-free).
type Manager struct {
	mu      sync.Mutex
	runtime wazero.Runtime
	hostAPI HostAPI
	modules map[string]api.Module
}

// NewManager creates a WASM Plugin Manager instance.
func NewManager(ctx context.Context, hostAPI HostAPI) (*Manager, error) {
	if hostAPI == nil {
		hostAPI = &DefaultHostAPI{}
	}

	r := wazero.NewRuntime(ctx)

	// Host functions
	_, err := r.NewHostModuleBuilder("env").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, ptr uint32, len uint32) {
			// WASM Host Log bridge
		}).Export("host_log").
		Instantiate(ctx)
	if err != nil {
		r.Close(ctx)
		return nil, fmt.Errorf("failed to build host module: %w", err)
	}

	return &Manager{
		runtime: r,
		hostAPI: hostAPI,
		modules: make(map[string]api.Module),
	}, nil
}

// LoadPlugin instantiates a compiled WASM plugin module bytes.
func (m *Manager) LoadPlugin(ctx context.Context, name string, wasmBytes []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	mod, err := m.runtime.Instantiate(ctx, wasmBytes)
	if err != nil {
		return fmt.Errorf("failed to instantiate WASM plugin %s: %w", name, err)
	}

	m.modules[name] = mod
	return nil
}

// Close gracefully closes the WASM runtime and all plugin modules.
func (m *Manager) Close(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.runtime != nil {
		return m.runtime.Close(ctx)
	}
	return nil
}
