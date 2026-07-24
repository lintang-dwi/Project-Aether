package tool_test

import (
	"context"
	"testing"
	"time"

	"aether/internal/tool"
)

func TestToolSandbox_RunCommand(t *testing.T) {
	tmpDir := t.TempDir()

	sandbox, err := tool.NewSandbox(tmpDir, 5*time.Second)
	if err != nil {
		t.Fatalf("failed to create sandbox: %v", err)
	}

	res, err := sandbox.RunCommand(context.Background(), "go", "version")
	if err != nil && res.ExitCode != 0 {
		t.Fatalf("command failed: %v, stderr: %s", err, res.Stderr)
	}

	if res.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", res.ExitCode)
	}
	if res.Stdout == "" {
		t.Error("expected non-empty stdout for 'go version'")
	}
}
