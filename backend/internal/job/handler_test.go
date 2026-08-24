package job

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/SsagarikaR/pipeline-processing/internal/models"
	"github.com/SsagarikaR/pipeline-processing/internal/pipeline"
)

func setupTestHandler() (*pipelineHandler, *mockJobStore, *mockResultStore, *mockErrorStore) {
	js := &mockJobStore{Jobs: make(map[uuid.UUID]models.Job)}
	rs := &mockResultStore{Results: make(map[uuid.UUID][]models.Result)}
	es := &mockErrorStore{Errors: make(map[uuid.UUID][]models.JobError)}
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
	if j.ID == uuid.Nil {
		t.Error("expected a generated job ID, got the zero UUID")
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
	id := uuid.New()
	js.Jobs[id] = models.Job{ID: id, Status: pipeline.StatusRunning}

	req := httptest.NewRequest(http.MethodGet, "/pipelines/"+id.String(), nil)
	req.SetPathValue("id", id.String())
	w := httptest.NewRecorder()

	h.GetJob(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestCancelJob_Success(t *testing.T) {
	h, js, _, _ := setupTestHandler()
	id := uuid.New()
	js.Jobs[id] = models.Job{ID: id, Status: pipeline.StatusRunning}

	_, cancel := context.WithCancel(context.Background())
	h.service.mu.Lock()
	h.service.cancelFuncs[id] = cancel
	h.service.mu.Unlock()

	req := httptest.NewRequest(http.MethodPatch, "/pipelines/"+id.String()+"/cancel", nil)
	req.SetPathValue("id", id.String())
	w := httptest.NewRecorder()

	h.CancelJob(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
}

func TestGetProgress_Success(t *testing.T) {
	h, js, _, _ := setupTestHandler()
	id := uuid.New()
	js.Jobs[id] = models.Job{ID: id, Status: pipeline.StatusRunning, TotalRecords: 100}

	h.service.mu.Lock()
	h.service.trackers[id] = pipeline.NewTracker(func(e pipeline.ProcessError) {})
	h.service.mu.Unlock()

	req := httptest.NewRequest(http.MethodGet, "/pipelines/"+id.String()+"/progress", nil)
	req.SetPathValue("id", id.String())
	w := httptest.NewRecorder()

	h.GetProgress(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestGetAllJobs_Success(t *testing.T) {
	h, js, _, _ := setupTestHandler()
	id := uuid.New()
	js.Jobs[id] = models.Job{ID: id}

	req := httptest.NewRequest(http.MethodGet, "/pipelines", nil)
	w := httptest.NewRecorder()

	h.GetAllJobs(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestDeleteJobs_Success(t *testing.T) {
	h, js, _, _ := setupTestHandler()
	id := uuid.New()
	js.Jobs[id] = models.Job{ID: id}

	req := httptest.NewRequest(http.MethodDelete, "/pipelines/"+id.String(), nil)
	req.SetPathValue("id", id.String())
	w := httptest.NewRecorder()

	h.DeleteJobs(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
}

func TestGetResults_Success(t *testing.T) {
	h, _, rs, _ := setupTestHandler()
	id := uuid.New()
	rs.Results[id] = []models.Result{{JobID: id}}

	req := httptest.NewRequest(http.MethodGet, "/pipelines/"+id.String()+"/results", nil)
	req.SetPathValue("id", id.String())
	w := httptest.NewRecorder()

	h.GetResults(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestGetErrors_Success(t *testing.T) {
	h, _, _, es := setupTestHandler()
	id := uuid.New()
	es.Errors[id] = []models.JobError{{JobID: id}}

	req := httptest.NewRequest(http.MethodGet, "/pipelines/"+id.String()+"/errors", nil)
	req.SetPathValue("id", id.String())
	w := httptest.NewRecorder()

	h.GetErrors(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
