package job

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	_ "github.com/SsagarikaR/pipeline-processing/internal/models"
	"github.com/SsagarikaR/pipeline-processing/internal/pipeline"
)

// NewPipelineHandler wires a JobService into an HTTP handler for the
// pipeline endpoints (create, list, progress, cancel, delete, etc).
func NewPipelineHandler(service *JobService) *pipelineHandler {
	return &pipelineHandler{service: service}
}

// CreateJob godoc
// @Summary Create a pipeline job
// @Description Create a new pipeline job with specified input, transform, aggregate, and export stages
// @Accept json
// @Produce json
// @Param job body pipeline.JobSpec true "Job Specification"
// @Success 201 {object} models.Job
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/pipelines [post]
func (h *pipelineHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var spec pipeline.JobSpec

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	if err := json.NewDecoder(r.Body).Decode(&spec); err != nil {
		http.Error(w, "failed to parse request: body too large or invalid json", http.StatusBadRequest)
		return
	}

	newJob, err := h.service.CreateJob(r.Context(), spec)
	if errors.Is(err, ErrInvalidJobType) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		slog.Error("failed to create job", "err", err)
		http.Error(w, "failed to create job", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newJob)
}

// GetJob godoc
// @Summary Get a pipeline job
// @Description Get details of a specific pipeline job by ID
// @Produce json
// @Param id path string true "Job ID (UUID)"
// @Success 200 {object} models.Job
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /pipelines/{id} [get]
func (h *pipelineHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	j, err := h.service.GetJob(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("failed to get job", "err", err)
		http.Error(w, "failed to get job", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(j)
}

// GetAllJobs godoc
// @Summary Get all pipeline jobs
// @Description Get a list of all pipeline jobs
// @Produce json
// @Success 200 {array} models.Job
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/pipelines [get]
func (h *pipelineHandler) GetAllJobs(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePagination(r, 50)
	jobs, err := h.service.GetAllJobs(r.Context(), limit, offset)
	if err != nil {
		slog.Error("failed to get all jobs", "err", err)
		http.Error(w, "failed to get all jobs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(jobs)
}

// DeleteJobs godoc
// @Summary Delete a pipeline job
// @Description Delete a specific pipeline job and its artifacts by ID
// @Param id path string true "Job ID (UUID)"
// @Success 204 "No Content"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1pipelines/{id} [delete]
func (h *pipelineHandler) DeleteJobs(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.DeleteJob(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("failed to delete job", "err", err)
		http.Error(w, "failed to delete job", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CancelJob handles PATCH /api/v1/pipelines/:id/cancel
// @Summary Cancel a pipeline job
// @Description Cancel a running pipeline job by ID
// @Param id path string true "Job ID (UUID)"
// @Success 204 "No Content"
// @Failure 400 {string} string "Bad Request"
// @Failure 409 {string} string "Conflict"
// @Router /api/v1/pipelines/{id}/cancel [patch]
func (h *pipelineHandler) CancelJob(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CancelJob(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusConflict) // job already finished / not running
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetProgress handles GET /api/v1/pipelines/:id/progress
// @Summary Get pipeline progress
// @Description Get progress and metrics of a specific pipeline job by ID
// @Produce json
// @Param id path string true "Job ID (UUID)"
// @Success 200 {object} interface{}
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/pipelines/{id}/progress [get]
func (h *pipelineHandler) GetProgress(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	j, processed, errCount, latencies, err := h.service.GetProgress(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("failed to get progress", "err", err)
		http.Error(w, "failed to get progress", http.StatusInternalServerError)
		return
	}

	percent := 0.0
	if j.TotalRecords > 0 {
		percent = float64(processed) / float64(j.TotalRecords) * 100
	}
	if j.Status == "completed" {
		percent = 100.0
	}

	var rate float64
	if j.StartedAt != nil {
		end := time.Now()
		if j.CompletedAt != nil {
			end = *j.CompletedAt // freeze the rate once the job is done, don't keep diluting it against wall-clock time
		}
		elapsed := end.Sub(*j.StartedAt).Seconds()
		if elapsed > 0 {
			rate = float64(processed) / elapsed
		}
	}

	resp := struct {
		JobID           uuid.UUID         `json:"jobId"`
		Status          string            `json:"status"`
		Processed       int64             `json:"processed"`
		ErrorCount      int64             `json:"errorCount"`
		PercentComplete float64           `json:"percentComplete"`
		RecordsPerSec   float64           `json:"recordsPerSec"`
		StageLatencies  map[string]string `json:"stageLatencies,omitempty"`
		StartedAt       any               `json:"startedAt"`
		CompletedAt     any               `json:"completedAt"`
		ExportUrl       *string           `json:"exportUrl"`
	}{j.ID, j.Status, processed, errCount, percent, rate, latencies, j.StartedAt, j.CompletedAt, j.ExportURL}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// GetResults handles GET /api/v1/pipelines/:id/results
// @Summary Get pipeline results
// @Description Retrieve results for a specific pipeline job by ID
// @Produce json
// @Param id path string true "Job ID (UUID)"
// @Success 200 {array} models.Result
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/pipelines/{id}/results [get]
func (h *pipelineHandler) GetResults(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	limit, offset := parsePagination(r, 50)
	results, err := h.service.GetResults(r.Context(), id, limit, offset)
	if err != nil {
		slog.Error("failed to get results", "err", err)
		http.Error(w, "failed to get results", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(results)
}

// GetErrors handles GET /api/v1/pipelines/:id/errors
// @Summary Get pipeline errors
// @Description Retrieve job errors and failed records for a specific pipeline job by ID
// @Produce json
// @Param id path string true "Job ID (UUID)"
// @Success 200 {array} models.JobError
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/pipelines/{id}/errors [get]
func (h *pipelineHandler) GetErrors(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	limit, offset := parsePagination(r, 50)
	errs, err := h.service.GetErrors(r.Context(), id, limit, offset)
	if err != nil {
		slog.Error("failed to get errors", "err", err)
		http.Error(w, "failed to get errors", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(errs)
}

// parseID reads the "id" path value from the request and parses it as a
// UUID, returning an error if it's missing or not a valid one.
func parseID(r *http.Request) (uuid.UUID, error) {
	idStr := r.PathValue("id")
	if idStr == "" {
		return uuid.UUID{}, errors.New("missing job ID")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.UUID{}, errors.New("invalid job ID")
	}
	return id, nil
}

// parsePagination extracts limit and offset from the query string,
// providing a safe default limit.
func parsePagination(r *http.Request, defaultLimit int) (int, int) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := defaultLimit
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	offset := 0
	if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
		offset = o
	}
	return limit, offset
}
