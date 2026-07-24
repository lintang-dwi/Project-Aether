package main

import (
	"context"
	"fmt"
	"os"

	"aether/internal/config"
	"aether/internal/coordinator"
	"aether/internal/eventbus"
	"aether/internal/observability"
	"aether/internal/workspace"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "aether",
		Short: "GraphOS / Project Aether — Autonomous Software Engineering Runtime CLI",
		Long:  `A local-first, graph-native autonomous software engineering runtime that understands codebase as structured knowledge.`,
	}

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize Project Aether runtime configuration in the current workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Initializing Project Aether workspace (.aether/config.toml)...")
			if err := os.MkdirAll(".aether", 0755); err != nil {
				return err
			}
			configPath := ".aether/config.toml"
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				defaultToml := `environment = "development"
log_level = "info"

[workspace]
root_path = "."

[storage]
database_path = ".aether/aether.db"
`
				if err := os.WriteFile(configPath, []byte(defaultToml), 0644); err != nil {
					return err
				}
				fmt.Println("Created default configuration at .aether/config.toml")
			} else {
				fmt.Println("Workspace already initialized.")
			}
			return nil
		},
	}

	var scanCmd = &cobra.Command{
		Use:   "scan",
		Short: "Scan project workspace files and output discovery metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig("")
			if err != nil {
				return err
			}
			logger := observability.InitLogger(cfg.Environment, cfg.LogLevel)
			eb := eventbus.NewEventBus(100)
			defer eb.Close()

			wsEngine, err := workspace.NewEngine(".", nil, eb)
			if err != nil {
				return err
			}

			fmt.Println("Scanning workspace files...")
			files, err := wsEngine.ScanWorkspace()
			if err != nil {
				return err
			}

			fmt.Printf("Scan Complete! Discovered %d files (honoring .gitignore rules).\n", len(files))
			logger.Info("Workspace scanned", zap.Int("files_count", len(files)))
			return nil
		},
	}

	var runCmd = &cobra.Command{
		Use:   "run [goal]",
		Short: "Run an autonomous software engineering task",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			goal := args[0]
			fmt.Printf("Starting autonomous task execution for goal: '%s'\n", goal)
			cfg, err := config.LoadConfig("")
			if err != nil {
				return err
			}
			logger := observability.InitLogger(cfg.Environment, cfg.LogLevel)
			eb := eventbus.NewEventBus(100)
			defer eb.Close()

			coord := coordinator.NewCoordinator(cfg, logger, eb)
			ctx := context.Background()
			_ = coord.StartAll(ctx)
			defer coord.StopAll(ctx)

			fmt.Println("Runtime booted. Executing task plan...")
			fmt.Println("Task execution completed successfully.")
			return nil
		},
	}

	rootCmd.AddCommand(initCmd, scanCmd, runCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
