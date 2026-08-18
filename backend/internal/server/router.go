package server

import (
	"database/sql"
	"net/http"

	"github.com/SsagarikaR/pipeline-processing/internal/job"
)

func mapRoutes(pool *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	// Initialize Job Dependencies
	jobStore := job.NewJobStore(pool)
	jobService := job.NewJobService(jobStore)
	ph := job.NewPipelineHandler(jobService)

	// Define Job Endpoints
	mux.HandleFunc("POST /pipelines", ph.CreateJob)
	mux.HandleFunc("GET /pipelines/{id}", ph.GetJob)
	mux.HandleFunc("GET /pipelines", ph.GetAllJobs)
	mux.HandleFunc("DELETE /pipelines/{id}", ph.DeleteJobs)

	
	return mux
}
