package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/SsagarikaR/pipeline-processing/internal/config"
	"github.com/SsagarikaR/pipeline-processing/internal/db"
	"github.com/SsagarikaR/pipeline-processing/pkg/logger"
	"github.com/SsagarikaR/pipeline-processing/internal/server"
)

// @title Pipeline Processing API
// @version 1.0
// @description This is a data processing pipeline API.
// @host localhost:8080
// @BasePath /
func main() {
	if err := run(); err != nil {
		slog.Error("startup failed", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.LoadConfig()
	slog.SetDefault(logger.New(cfg.LogLevel, cfg.LogFormat))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.New(ctx, cfg.DB)
	if err != nil {
		return fmt.Errorf("db connection failed: %w", err)
	}
	defer func() {
		if err := pool.Close(); err != nil {
			slog.Error("db close failed", "err", err)
		}
	}()

	srv := server.New(cfg, pool)

	go func() {
		slog.Info("server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutdown signal received")

	server.Shutdown(context.Background(), srv)
	slog.Info("server stopped")
	return nil
}
