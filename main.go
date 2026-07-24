package main

import (
	"context"
	"fmt"
	"os"

	"aether/internal/config"
	"aether/internal/coordinator"
	"aether/internal/eventbus"
	"aether/internal/observability"
	"aether/internal/ui"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"go.uber.org/zap"
)

const (
	Version = "0.1.0"
	AppName = "Project Aether"
)

func main() {
	// 1. Load Configuration & Validate
	cfg, err := config.LoadConfig("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// 2. Initialize Uber Zap Logger
	logger := observability.InitLogger(cfg.Environment, cfg.LogLevel)
	defer logger.Sync()

	logger.Info("Booting "+AppName+"...",
		zap.String("version", Version),
		zap.String("environment", cfg.Environment),
	)

	// 3. Initialize EventBus
	eb := eventbus.NewEventBus(cfg.EventBus.BufferSize)
	defer eb.Close()

	// 4. Initialize Coordinator
	coord := coordinator.NewCoordinator(cfg, logger, eb)
	ctx := context.Background()
	_ = coord.StartAll(ctx)

	// 5. Initialize Wails App Bridge
	appUI := ui.NewApp(coord)

	// 6. Run Wails Application
	err = wails.Run(&options.App{
		Title:  AppName + " — GraphOS Desktop",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: Assets,
		},
		BackgroundColour: &options.RGBA{R: 7, G: 8, B: 12, A: 255},
		OnStartup:        appUI.Startup,
		Bind: []interface{}{
			appUI,
		},
	})

	if err != nil {
		logger.Error("Wails application error", zap.Error(err))
		_ = coord.StopAll(ctx)
		os.Exit(1)
	}

	_ = coord.StopAll(ctx)
}
