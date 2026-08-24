package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/imabg/onboard/internal/config"
	"github.com/imabg/onboard/internal/db"
	"github.com/imabg/onboard/internal/logger"
	"github.com/imabg/onboard/internal/server"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	if err := logger.Init(cfg.Log); err != nil {
		panic(err)
	}
	defer logger.Sync()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.NewPool(ctx, cfg.Database)
	if err != nil {
		logger.L().Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	srv := server.New(cfg.Server, pool)

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.L().Error("graceful shutdown failed", zap.Error(err))
		}
	case err := <-errCh:
		if err != nil {
			logger.L().Fatal("http server error", zap.Error(err))
		}
	}

	logger.L().Info("server stopped")
}
