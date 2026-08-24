package logger

import (
	"fmt"
	"sync"

	"github.com/imabg/onboard/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	mu  sync.RWMutex
	log = zap.NewNop()
)

func Init(cfg config.LogConfig) error {
	l, err := build(cfg)
	if err != nil {
		return err
	}
	Replace(l)
	return nil
}

func L() *zap.Logger {
	mu.RLock()
	defer mu.RUnlock()
	return log
}

func Sync() {
	_ = L().Sync()
}

func Replace(l *zap.Logger) {
	if l == nil {
		l = zap.NewNop()
	}
	mu.Lock()
	log = l
	mu.Unlock()
	zap.ReplaceGlobals(l)
}

func build(cfg config.LogConfig) (*zap.Logger, error) {
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("parse log level %q: %w", cfg.Level, err)
	}

	zcfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       true,
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    "func",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	l, err := zcfg.Build(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		return nil, fmt.Errorf("build logger: %w", err)
	}
	return l, nil
}
