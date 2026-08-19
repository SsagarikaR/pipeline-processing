package job

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	_"github.com/SsagarikaR/pipeline-processing/internal/models"
)

type pipelineHandler struct {
	service *JobService
}

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
// @Router /pipelines [post]
func (h *pipelineHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var spec pipeline.JobSpec

	if err := json.NewDecoder(r.Body).Decode(&spec); err != nil {
		http.Error(w, "failed to parse request", http.StatusBadRequest)
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
	json.NewEncoder(w).Encode(newJob)
}

// GetJob godoc
// @Summary Get a pipeline job
// @Description Get details of a specific pipeline job by ID
// @Produce json
// @Param id path int true "Job ID"
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
	json.NewEncoder(w).Encode(j)
}

// GetAllJobs godoc
// @Summary Get all pipeline jobs
// @Description Get a list of all pipeline jobs
// @Produce json
// @Success 200 {array} models.Job
// @Failure 500 {string} string "Internal Server Error"
// @Router /pipelines [get]
func (h *pipelineHandler) GetAllJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.service.GetAllJobs(r.Context())
	if err != nil {
		slog.Error("failed to get all jobs", "err", err)
		http.Error(w, "failed to get all jobs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

// DeleteJobs godoc
// @Summary Delete a pipeline job
// @Description Delete a specific pipeline job and its artifacts by ID
// @Param id path int true "Job ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /pipelines/{id} [delete]
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

func parseID(r *http.Request) (int, error) {
	idStr := r.PathValue("id")
	if idStr == "" {
		return 0, errors.New("missing job ID")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("invalid job ID")
	}
	return id, nil
}