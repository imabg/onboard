package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/imabg/onboard/internal/config"
	"github.com/imabg/onboard/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Server struct {
	http *http.Server
	pool *pgxpool.Pool
}

func New(cfg config.ServerConfig, pool *pgxpool.Pool) *Server {
	s := &Server{pool: pool}
	s.http = &http.Server{
		Addr:         cfg.Addr(),
		Handler:      s.routes(),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	return s
}

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /ready", s.handleReady)
	mux.HandleFunc("GET /api/v1/", s.handleRoot)
	mux.HandleFunc("GET /api/v1", s.handleRoot)
	return loggingMiddleware(mux)
}

func (s *Server) Start() error {
	logger.L().Info("http server listening", zap.String("addr", s.http.Addr))
	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	logger.L().Info("http server shutting down")
	return s.http.Shutdown(ctx)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := s.pool.Ping(ctx); err != nil {
		logger.L().Error("readiness check failed", zap.Error(err))
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"status": "unavailable",
			"error":  "database unreachable",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"service": "onboard",
		"message": "HR portal API",
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
