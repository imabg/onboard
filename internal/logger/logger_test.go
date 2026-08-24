package logger

import (
	"testing"

	"github.com/imabg/onboard/internal/config"
)

func TestNewProduction(t *testing.T) {
	log, err := New(config.LogConfig{Level: "info", Encoding: "json", Development: false})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if log == nil {
		t.Fatal("New() returned nil logger")
	}
	_ = log.Sync()
}

func TestNewInvalidLevel(t *testing.T) {
	_, err := New(config.LogConfig{Level: "loud"})
	if err == nil {
		t.Fatal("New() expected error for invalid level")
	}
}
