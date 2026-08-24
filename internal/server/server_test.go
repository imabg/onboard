package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imabg/onboard/internal/config"
	"go.uber.org/zap"
)

func TestHandleHealth(t *testing.T) {
	s := New(config.ServerConfig{Host: "127.0.0.1", Port: 8080}, zap.NewNop(), nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	s.routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("status = %q, want ok", body["status"])
	}
}

func TestHandleRoot(t *testing.T) {
	s := New(config.ServerConfig{Host: "127.0.0.1", Port: 8080}, zap.NewNop(), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/", nil)
	rec := httptest.NewRecorder()
	s.routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["service"] != "onboard" {
		t.Errorf("service = %q, want onboard", body["service"])
	}
}

func TestUnknownRoute(t *testing.T) {
	s := New(config.ServerConfig{Host: "127.0.0.1", Port: 8080}, zap.NewNop(), nil)

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()
	s.routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}
