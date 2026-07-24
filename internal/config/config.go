package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/pelletier/go-toml/v2"
)

// Config holds the complete configuration for Project Aether.
type Config struct {
	Environment string            `toml:"environment" validate:"required"`
	LogLevel    string            `toml:"log_level" validate:"required"`
	Workspace   WorkspaceConfig   `toml:"workspace"`
	Storage     StorageConfig     `toml:"storage"`
	EventBus    EventBusConfig    `toml:"eventbus"`
	Custom      map[string]string `toml:"custom"`
}

type WorkspaceConfig struct {
	RootPath       string   `toml:"root_path" validate:"required"`
	IgnorePatterns []string `toml:"ignore_patterns"`
}

type StorageConfig struct {
	DatabasePath string `toml:"database_path" validate:"required"`
	InMemory     bool   `toml:"in_memory"`
}

type EventBusConfig struct {
	BufferSize int `toml:"buffer_size" validate:"gt=0"`
}

// DefaultConfig returns the sensible default configuration.
func DefaultConfig() Config {
	return Config{
		Environment: "development",
		LogLevel:    "info",
		Workspace: WorkspaceConfig{
			RootPath:       ".",
			IgnorePatterns: []string{".git", "node_modules", "vendor", "bin", "dist"},
		},
		Storage: StorageConfig{
			DatabasePath: ".aether/aether.db",
			InMemory:     false,
		},
		EventBus: EventBusConfig{
			BufferSize: 1024,
		},
		Custom: make(map[string]string),
	}
}

// Validate validates the config struct using go-playground/validator.
func (c *Config) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

// LoadConfig loads configuration with layered precedence:
// Defaults -> File (.aether/config.toml) -> Environment Overrides -> Validation
func LoadConfig(configPath string) (Config, error) {
	cfg := DefaultConfig()

	if configPath == "" {
		configPath = filepath.Join(".aether", "config.toml")
	}

	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return cfg, fmt.Errorf("failed to read config file %s: %w", configPath, err)
		}
		if err := toml.Unmarshal(data, &cfg); err != nil {
			return cfg, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return cfg, fmt.Errorf("error checking config file %s: %w", configPath, err)
	}

	// Environment variable overrides
	if env := os.Getenv("AETHER_ENV"); env != "" {
		cfg.Environment = env
	}
	if level := os.Getenv("AETHER_LOG_LEVEL"); level != "" {
		cfg.LogLevel = level
	}
	if dbPath := os.Getenv("AETHER_DB_PATH"); dbPath != "" {
		cfg.Storage.DatabasePath = dbPath
	}

	if err := cfg.Validate(); err != nil {
		return cfg, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}
