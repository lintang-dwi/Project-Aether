package security_test

import (
	"testing"

	"aether/internal/security"
	"go.uber.org/zap"
)

func TestSecurityBoundary_Permissions(t *testing.T) {
	logger := zap.NewNop()
	sec := security.NewBoundary(true, false, logger) // Allow write, disallow execute

	// Read -> Should pass
	if err := sec.CheckPermission(security.PermRead, "main.go"); err != nil {
		t.Errorf("expected READ allowed, got %v", err)
	}

	// Write -> Should pass
	if err := sec.CheckPermission(security.PermWrite, "main.go"); err != nil {
		t.Errorf("expected WRITE allowed, got %v", err)
	}

	// Execute -> Should fail
	if err := sec.CheckPermission(security.PermExecute, "rm -rf /"); err == nil {
		t.Error("expected EXECUTE to be blocked by security boundary")
	}

	logs := sec.GetAuditLogs()
	if len(logs) != 3 {
		t.Errorf("expected 3 audit logs, got %d", len(logs))
	}
}
