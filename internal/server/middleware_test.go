package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/imabg/onboard/internal/config"
	"github.com/imabg/onboard/internal/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggingMiddlewareRecordsStatus(t *testing.T) {
	core, logs := observer.New(zapcore.InfoLevel)
	logger.Replace(zap.New(core))
	t.Cleanup(func() { logger.Replace(zap.NewNop()) })

	s := New(config.ServerConfig{Host: "127.0.0.1", Port: 8080}, nil)

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
	if fields["payload"] != "" {
		t.Errorf("payload = %v, want empty string", fields["payload"])
	}
}

func TestLoggingMiddlewareRecordsJSONPayload(t *testing.T) {
	core, logs := observer.New(zapcore.InfoLevel)
	logger.Replace(zap.New(core))
	t.Cleanup(func() { logger.Replace(zap.NewNop()) })

	s := New(config.ServerConfig{Host: "127.0.0.1", Port: 8080}, nil)

	body := `{"email":"abhay@example.com","name":"Abhay"}`
	req := httptest.NewRequest(http.MethodGet, "/health?ref=signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.routes().ServeHTTP(rec, req)

	entries := logs.FilterMessage("request").All()
	if len(entries) != 1 {
		t.Fatalf("got %d request logs, want 1", len(entries))
	}

	fields := entries[0].ContextMap()
	payload, ok := fields["payload"].(map[string]any)
	if !ok {
		t.Fatalf("payload type = %T, want map[string]any", fields["payload"])
	}
	if payload["email"] != "abhay@example.com" {
		t.Errorf("payload.email = %v", payload["email"])
	}
	if fields["query"] != "ref=signup" {
		t.Errorf("query = %v, want ref=signup", fields["query"])
	}
	if fields["content_type"] != "application/json" {
		t.Errorf("content_type = %v", fields["content_type"])
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
