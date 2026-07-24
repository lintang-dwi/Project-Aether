package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"aether/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg.Environment != "development" {
		t.Errorf("expected default environment 'development', got '%s'", cfg.Environment)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("expected default log level 'info', got '%s'", cfg.LogLevel)
	}
	if cfg.EventBus.BufferSize != 1024 {
		t.Errorf("expected buffer size 1024, got %d", cfg.EventBus.BufferSize)
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected valid config, got error: %v", err)
	}
}

func TestLoadConfig_FileAndEnv(t *testing.T) {
	tmpDir := t.TempDir()
	cfgFile := filepath.Join(tmpDir, "config.toml")
	tomlData := `
environment = "production"
log_level = "debug"

[workspace]
root_path = "/tmp/project"

[storage]
database_path = "/tmp/aether.db"

[eventbus]
buffer_size = 512
`
	if err := os.WriteFile(cfgFile, []byte(tomlData), 0644); err != nil {
		t.Fatalf("failed to write test config file: %v", err)
	}

	// Test loading from file
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		t.Fatalf("unexpected error loading config: %v", err)
	}
	if cfg.Environment != "production" {
		t.Errorf("expected 'production', got '%s'", cfg.Environment)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("expected 'debug', got '%s'", cfg.LogLevel)
	}
	if cfg.Workspace.RootPath != "/tmp/project" {
		t.Errorf("expected '/tmp/project', got '%s'", cfg.Workspace.RootPath)
	}

	// Test Env override
	t.Setenv("AETHER_ENV", "staging")
	t.Setenv("AETHER_LOG_LEVEL", "warn")

	cfgEnv, err := config.LoadConfig(cfgFile)
	if err != nil {
		t.Fatalf("unexpected error loading config with env: %v", err)
	}
	if cfgEnv.Environment != "staging" {
		t.Errorf("expected env override 'staging', got '%s'", cfgEnv.Environment)
	}
	if cfgEnv.LogLevel != "warn" {
		t.Errorf("expected env override 'warn', got '%s'", cfgEnv.LogLevel)
	}
}
