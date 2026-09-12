package logger

import (
	"fmt"

	"go.uber.org/zap"
)

var Log *zap.Logger = zap.NewNop()

func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return fmt.Errorf("failed to parse log lvl: %w", err)
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl

	logger, err := cfg.Build()
	if err != nil {
		return fmt.Errorf("failed to build logger: %w", err)
	}

	Log = logger
	return nil
}
