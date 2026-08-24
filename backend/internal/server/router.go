package server

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/SsagarikaR/pipeline-processing/docs"
	"github.com/SsagarikaR/pipeline-processing/internal/job"
)

// mapRoutes builds the job dependencies (store, service, handler) and
// registers every HTTP endpoint the API exposes, plus the swagger UI.
func mapRoutes(ctx context.Context, pool *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	p := func(path string) string {
		return "/api/v1" + path
	}
	// Swagger endpoint
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	// Initialize Job Dependencies
	jobStore := job.NewJobStore(pool)
	resultStore := job.NewResultStore(pool)
	errorStore := job.NewErrorStore(pool)
	jobService := job.NewJobService(jobStore, resultStore, errorStore)

	// Best-effort cleanup: don't block startup on it, but do log if it fails.
	if err := jobService.RecoverStuckJobs(ctx); err != nil {
		slog.Error("failed to recover stuck jobs", "err", err)
	}

	ph := job.NewPipelineHandler(jobService)

	// Define Job Endpoints
	mux.HandleFunc("POST "+p("/pipelines"), ph.CreateJob)
	mux.HandleFunc("GET "+p("/pipelines/{id}"), ph.GetJob)
	mux.HandleFunc("GET "+p("/pipelines"), ph.GetAllJobs)
	mux.HandleFunc("GET "+p("/pipelines/{id}/progress"), ph.GetProgress)
	mux.HandleFunc("GET "+p("/pipelines/{id}/results"), ph.GetResults)
	mux.HandleFunc("GET "+p("/pipelines/{id}/errors"), ph.GetErrors)
	mux.HandleFunc("PATCH "+p("/pipelines/{id}/cancel"), ph.CancelJob)
	mux.HandleFunc("DELETE "+p("/pipelines/{id}"), ph.DeleteJobs)

	return mux
}
