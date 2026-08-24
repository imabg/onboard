package logger

import (
	"testing"

	"github.com/imabg/onboard/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestInit(t *testing.T) {
	t.Cleanup(func() { Replace(zap.NewNop()) })

	if err := Init(config.LogConfig{Level: "info"}); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if L() == nil {
		t.Fatal("L() returned nil")
	}
	Sync()
}

func TestInitInvalidLevel(t *testing.T) {
	if err := Init(config.LogConfig{Level: "loud"}); err == nil {
		t.Fatal("Init() expected error for invalid level")
	}
}

func TestReplace(t *testing.T) {
	core, logs := observer.New(zapcore.InfoLevel)
	Replace(zap.New(core))
	t.Cleanup(func() { Replace(zap.NewNop()) })

	L().Info("hello")
	if logs.Len() != 1 {
		t.Fatalf("got %d logs, want 1", logs.Len())
	}
	if logs.All()[0].Message != "hello" {
		t.Errorf("message = %q, want hello", logs.All()[0].Message)
	}
}
