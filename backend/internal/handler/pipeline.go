package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/SsagarikaR/pipeline-processing/internal/store"
)

type pipelineHandler struct{
	store store.JobStore
}

func NewPipelineHandler(store store.JobStore) *pipelineHandler {
	return &pipelineHandler{store: store}
}

func (h *pipelineHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var req struct{
		Type string `json:"type"`
		Data []int  `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "failed to parse request", http.StatusBadRequest)
		return
	}

	jobSpec, err := json.Marshal(req)
	if err != nil {
		http.Error(w, "failed to serialize job spec", http.StatusInternalServerError)
		return
	}

	job,err := h.store.CreateJob(r.Context(), jobSpec)
	if err != nil {
		slog.Error("failed to create job", "err", err)
		http.Error(w, "failed to create job", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(job)
}

func (h *pipelineHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	jobIDStr := r.PathValue("id")
	if jobIDStr == "" {
		http.Error(w, "missing job ID", http.StatusBadRequest)
		return
	}

	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		http.Error(w, "invalid job ID", http.StatusBadRequest)
		return
	}

	job,err := h.store.GetJob(r.Context(), jobID)
	if err != nil {
		slog.Error("failed to get job", "err", err)
		http.Error(w, "failed to get job", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func (h *pipelineHandler) GetAllJobs(w http.ResponseWriter, r *http.Request) {
	jobs,err := h.store.GetAllJobs(r.Context())
	if err != nil {
		slog.Error("failed to get all jobs", "err", err)
		http.Error(w, "failed to get all jobs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

func (h *pipelineHandler) DeleteJobs(w http.ResponseWriter, r *http.Request){
	jobID := r.PathValue("id")
	id, err := strconv.Atoi(jobID)
	if err != nil {
		http.Error(w, "invalid job id", http.StatusBadRequest)
		return
	}
	err = h.store.DeleteJobs(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows){
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