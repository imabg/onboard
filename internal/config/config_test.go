package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/testdb?sslmode=disable")
	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("LOG_LEVEL", "debug")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Server.Port != 9090 {
		t.Errorf("Server.Port = %d, want 9090", cfg.Server.Port)
	}
	if cfg.Database.URL != "postgres://user:pass@localhost:5432/testdb?sslmode=disable" {
		t.Errorf("Database.URL = %q", cfg.Database.URL)
	}
	if cfg.Log.Level != "debug" {
		t.Errorf("Log.Level = %q, want debug", cfg.Log.Level)
	}
	if cfg.Server.Addr() != "0.0.0.0:9090" {
		t.Errorf("Addr() = %q, want 0.0.0.0:9090", cfg.Server.Addr())
	}
	if cfg.Server.ReadTimeout != 15*time.Second {
		t.Errorf("ReadTimeout = %v, want 15s", cfg.Server.ReadTimeout)
	}
}

func TestLoadInvalidPort(t *testing.T) {
	t.Setenv("SERVER_PORT", "70000")
	t.Setenv("DATABASE_URL", "postgres://localhost/db")

	if _, err := Load(); err == nil {
		t.Fatal("Load() expected error for invalid port")
	}
}

func TestLoadFromYAML(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir("configs", 0o755); err != nil {
		t.Fatal(err)
	}

	yaml := []byte(`
server:
  host: "127.0.0.1"
  port: 3000
database:
  url: "postgres://yaml:yaml@db:5432/onboard"
log:
  level: "warn"
  encoding: "console"
  development: true
`)
	if err := os.WriteFile("configs/config.yaml", yaml, 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("Server.Host = %q, want 127.0.0.1", cfg.Server.Host)
	}
	if cfg.Server.Port != 3000 {
		t.Errorf("Server.Port = %d, want 3000", cfg.Server.Port)
	}
	if cfg.Database.URL != "postgres://yaml:yaml@db:5432/onboard" {
		t.Errorf("Database.URL = %q", cfg.Database.URL)
	}
	if cfg.Log.Level != "warn" {
		t.Errorf("Log.Level = %q, want warn", cfg.Log.Level)
	}
	if !cfg.Log.Development {
		t.Error("Log.Development = false, want true")
	}
}
