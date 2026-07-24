package observability_test

import (
	"testing"

	"aether/internal/observability"
	"go.uber.org/zap/zapcore"
)

func TestInitLogger_ZapDevelopment(t *testing.T) {
	logger := observability.InitLogger("development", "debug")
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}

	if !logger.Core().Enabled(zapcore.DebugLevel) {
		t.Error("expected debug level to be enabled")
	}
}

func TestInitLogger_ZapProduction(t *testing.T) {
	logger := observability.InitLogger("production", "info")
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}

	if logger.Core().Enabled(zapcore.DebugLevel) {
		t.Error("expected debug level to be disabled in info mode")
	}
}
