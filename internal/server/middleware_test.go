package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imabg/onboard/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggingMiddlewareRecordsStatus(t *testing.T) {
	core, logs := observer.New(zapcore.InfoLevel)
	s := New(config.ServerConfig{Host: "127.0.0.1", Port: 8080}, zap.New(core), nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	s.routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	entries := logs.FilterMessage("request").All()
	if len(entries) != 1 {
		t.Fatalf("got %d request logs, want 1", len(entries))
	}

	fields := entries[0].ContextMap()
	if fields["method"] != http.MethodGet {
		t.Errorf("method = %v, want GET", fields["method"])
	}
	if fields["path"] != "/health" {
		t.Errorf("path = %v, want /health", fields["path"])
	}
	if fields["status"] != int64(http.StatusOK) {
		t.Errorf("status = %v, want 200", fields["status"])
	}
}

func TestStatusWriterCapturesCode(t *testing.T) {
	rec := httptest.NewRecorder()
	w := &statusWriter{ResponseWriter: rec, status: http.StatusOK}
	w.WriteHeader(http.StatusCreated)

	if w.status != http.StatusCreated {
		t.Errorf("status = %d, want 201", w.status)
	}
	if rec.Code != http.StatusCreated {
		t.Errorf("recorder code = %d, want 201", rec.Code)
	}
}
