package server

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/SsagarikaR/pipeline-processing/internal/config"
)

func New(cfg *config.Config, pool *sql.DB) *http.Server {
	router := mapRoutes(pool)

	return &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router, 
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func Shutdown(ctx context.Context, srv *http.Server) {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "err", err)
	}
}
