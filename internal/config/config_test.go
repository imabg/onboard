package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

const validYAML = `
server:
  host: "127.0.0.1"
  port: 3000
  read_timeout: 15s
  write_timeout: 15s
  idle_timeout: 60s
  shutdown_timeout: 10s
database:
  url: "postgres://yaml:yaml@db:5432/onboard"
  max_conns: 10
  min_conns: 2
  max_conn_lifetime: 1h
  max_conn_idle_time: 30m
log:
  level: "warn"
`

func TestLoadFromYAML(t *testing.T) {
	chdirWithConfig(t, validYAML)

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
	if cfg.Server.Addr() != "127.0.0.1:3000" {
		t.Errorf("Addr() = %q, want 127.0.0.1:3000", cfg.Server.Addr())
	}
	if cfg.Server.ReadTimeout != 15*time.Second {
		t.Errorf("ReadTimeout = %v, want 15s", cfg.Server.ReadTimeout)
	}
	if cfg.Database.URL != "postgres://yaml:yaml@db:5432/onboard" {
		t.Errorf("Database.URL = %q", cfg.Database.URL)
	}
	if cfg.Database.MaxConns != 10 {
		t.Errorf("MaxConns = %d, want 10", cfg.Database.MaxConns)
	}
	if cfg.Log.Level != "warn" {
		t.Errorf("Log.Level = %q, want warn", cfg.Log.Level)
	}
}

func TestLoadMissingFile(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	if _, err := Load(); err == nil {
		t.Fatal("Load() expected error for missing config file")
	}
}

func TestLoadInvalidPort(t *testing.T) {
	chdirWithConfig(t, `
server:
  host: "127.0.0.1"
  port: 70000
database:
  url: "postgres://localhost/db"
log:
  level: "info"
`)

	if _, err := Load(); err == nil {
		t.Fatal("Load() expected error for invalid port")
	}
}

func TestLoadMissingDatabaseURL(t *testing.T) {
	chdirWithConfig(t, `
server:
  host: "127.0.0.1"
  port: 8080
database:
  url: ""
log:
  level: "info"
`)

	if _, err := Load(); err == nil {
		t.Fatal("Load() expected error for missing database.url")
	}
}

func chdirWithConfig(t *testing.T, contents string) {
	t.Helper()
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "configs"), 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "configs", "config.yaml")
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
	chdir(t, dir)
}

func chdir(t *testing.T, dir string) {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })
}
