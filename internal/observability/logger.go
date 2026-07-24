package observability

import (
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	globalLogger *zap.Logger
	once         sync.Once
)

// InitLogger configures and initializes the global Uber Zap logger according to MVP Stack recommendations.
func InitLogger(env, logLevel string) *zap.Logger {
	var config zap.Config

	if strings.ToLower(env) == "development" {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
	}

	level := parseLogLevel(logLevel)
	config.Level = zap.NewAtomicLevelAt(level)

	logger, err := config.Build()
	if err != nil {
		logger = zap.NewNop()
	}

	once.Do(func() {
		globalLogger = logger
	})

	return logger
}

// GetLogger returns the global Zap logger.
func GetLogger() *zap.Logger {
	if globalLogger == nil {
		return zap.L()
	}
	return globalLogger
}

func parseLogLevel(levelStr string) zapcore.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
