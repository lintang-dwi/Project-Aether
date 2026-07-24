package memory_test

import (
	"testing"

	"aether/internal/memory"
)

func TestMemoryEngine_StoreAndRetrieve(t *testing.T) {
	mem := memory.NewEngine()

	mem.Store("decisions", "arch_style", "microkernel")

	rec, ok := mem.Retrieve("arch_style")
	if !ok {
		t.Fatal("expected record to exist")
	}
	if rec.Value != "microkernel" {
		t.Errorf("expected 'microkernel', got '%s'", rec.Value)
	}
}
