package job

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SsagarikaR/pipeline-processing/internal/models"
	"github.com/SsagarikaR/pipeline-processing/internal/pipeline"
)

func setupTestHandler() (*pipelineHandler, *mockJobStore, *mockResultStore, *mockErrorStore) {
	js := &mockJobStore{Jobs: make(map[int]models.Job)}
	rs := &mockResultStore{Results: make(map[int][]models.Result)}
	es := &mockErrorStore{Errors: make(map[int][]models.JobError)}
	svc := NewJobService(js, rs, es)
	return NewPipelineHandler(svc), js, rs, es
}

func TestCreateJob_Success(t *testing.T) {
	h, _, _, _ := setupTestHandler()

	spec := pipeline.JobSpec{
		Sources: []pipeline.SourceConfig{{Type: "csv", Path: "test.csv"}},
	}
	body, _ := json.Marshal(spec)

	req := httptest.NewRequest(http.MethodPost, "/pipelines", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateJob(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var j models.Job
	if err := json.NewDecoder(w.Body).Decode(&j); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if j.ID != 1 {
		t.Errorf("expected job ID 1, got %d", j.ID)
	}
}

func TestCreateJob_InvalidSpec(t *testing.T) {
	h, _, _, _ := setupTestHandler()

	spec := pipeline.JobSpec{} 
	body, _ := json.Marshal(spec)

	req := httptest.NewRequest(http.MethodPost, "/pipelines", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateJob(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestGetJob_Success(t *testing.T) {
	h, js, _, _ := setupTestHandler()
	js.Jobs[1] = models.Job{ID: 1, Status: pipeline.StatusRunning}

	req := httptest.NewRequest(http.MethodGet, "/pipelines/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.GetJob(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestCancelJob_Success(t *testing.T) {
	h, js, _, _ := setupTestHandler()
	js.Jobs[1] = models.Job{ID: 1, Status: pipeline.StatusRunning}

	_, cancel := context.WithCancel(context.Background())
	h.service.mu.Lock()
	h.service.cancelFuncs[1] = cancel
	h.service.mu.Unlock()

	req := httptest.NewRequest(http.MethodPatch, "/pipelines/1/cancel", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.CancelJob(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
}

func TestGetProgress_Success(t *testing.T) {
	h, js, _, _ := setupTestHandler()
	js.Jobs[1] = models.Job{ID: 1, Status: pipeline.StatusRunning, TotalRecords: 100}

	h.service.mu.Lock()
	h.service.trackers[1] = pipeline.NewTracker(func(e pipeline.ProcessError) {})
	h.service.mu.Unlock()

	req := httptest.NewRequest(http.MethodGet, "/pipelines/1/progress", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.GetProgress(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestGetAllJobs_Success(t *testing.T) {
	h, js, _, _ := setupTestHandler()
	js.Jobs[1] = models.Job{ID: 1}

	req := httptest.NewRequest(http.MethodGet, "/pipelines", nil)
	w := httptest.NewRecorder()

	h.GetAllJobs(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestDeleteJobs_Success(t *testing.T) {
	h, js, _, _ := setupTestHandler()
	js.Jobs[1] = models.Job{ID: 1}

	req := httptest.NewRequest(http.MethodDelete, "/pipelines/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.DeleteJobs(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
}

func TestGetResults_Success(t *testing.T) {
	h, _, rs, _ := setupTestHandler()
	rs.Results[1] = []models.Result{{JobID: 1}}

	req := httptest.NewRequest(http.MethodGet, "/pipelines/1/results", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.GetResults(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestGetErrors_Success(t *testing.T) {
	h, _, _, es := setupTestHandler()
	es.Errors[1] = []models.JobError{{JobID: 1}}

	req := httptest.NewRequest(http.MethodGet, "/pipelines/1/errors", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.GetErrors(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
