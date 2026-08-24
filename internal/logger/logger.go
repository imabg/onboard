package logger

import (
	"fmt"

	"github.com/imabg/onboard/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(cfg config.LogConfig) (*zap.Logger, error) {
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("parse log level %q: %w", cfg.Level, err)
	}

	var zcfg zap.Config
	if cfg.Development {
		zcfg = zap.NewDevelopmentConfig()
	} else {
		zcfg = zap.NewProductionConfig()
	}

	zcfg.Level = zap.NewAtomicLevelAt(level)
	if cfg.Encoding != "" {
		zcfg.Encoding = cfg.Encoding
	}
	zcfg.EncoderConfig.TimeKey = "ts"
	zcfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	log, err := zcfg.Build()
	if err != nil {
		return nil, fmt.Errorf("build logger: %w", err)
	}

	return log, nil
}
