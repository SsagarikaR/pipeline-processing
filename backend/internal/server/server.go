package server

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/SsagarikaR/pipeline-processing/internal/config"
	"github.com/SsagarikaR/pipeline-processing/internal/handler"
	"github.com/SsagarikaR/pipeline-processing/internal/store"
)

func New(cfg *config.Config, pool *sql.DB) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handler.Health(pool))

	jobStore := store.NewJobStore(pool)
	ph := handler.NewPipelineHandler(jobStore)
	mux.HandleFunc("POST /job", ph.CreateJob)
	mux.HandleFunc("GET /jobs/{id}", ph.GetJob)
	mux.HandleFunc("GET /jobs", ph.GetAllJobs)
	mux.HandleFunc("DELETE /jobs/{id}", ph.DeleteJobs)


	return &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
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