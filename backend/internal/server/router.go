package server

import (
	"database/sql"
	"net/http"

	_ "github.com/SsagarikaR/pipeline-processing/docs"
	"github.com/SsagarikaR/pipeline-processing/internal/job"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func mapRoutes(pool *sql.DB) *http.ServeMux {
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
