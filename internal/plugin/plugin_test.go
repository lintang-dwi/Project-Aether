package plugin_test

import (
	"context"
	"testing"

	"aether/internal/plugin"
)

func TestPluginManager_Initialization(t *testing.T) {
	ctx := context.Background()
	mgr, err := plugin.NewManager(ctx, nil)
	if err != nil {
		t.Fatalf("failed to create plugin manager: %v", err)
	}

	defer func() {
		if err := mgr.Close(ctx); err != nil {
			t.Errorf("error closing manager: %v", err)
		}
	}()
}
