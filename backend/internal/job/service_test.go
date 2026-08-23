package job

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/SsagarikaR/pipeline-processing/internal/models"
	"github.com/SsagarikaR/pipeline-processing/internal/pipeline"
)

type mockJobStore struct {
	Jobs map[int]models.Job
}

func (m *mockJobStore) CreateJob(ctx context.Context, spec json.RawMessage) (models.Job, error) {
	j := models.Job{ID: len(m.Jobs) + 1, Status: pipeline.StatusPending}
	m.Jobs[j.ID] = j
	return j, nil
}

func (m *mockJobStore) UpdateStatus(ctx context.Context, id int, status string) error {
	j, ok := m.Jobs[id]
	if !ok {
		return sql.ErrNoRows
	}
	j.Status = status
	m.Jobs[id] = j
	return nil
}

func (m *mockJobStore) UpdateStatusAndMetrics(ctx context.Context, jobID int, status string, processed int64, errors int64) error {
	j, ok := m.Jobs[jobID]
	if !ok {
		return sql.ErrNoRows
	}
	j.Status = status
	j.ProcessedRecords = int(processed)
	j.ErrorCount = int(errors)
	m.Jobs[jobID] = j
	return nil
}

func (m *mockJobStore) UpdateExportURL(ctx context.Context, jobID int, url string) error {
	j, ok := m.Jobs[jobID]
	if !ok {
		return sql.ErrNoRows
	}
	u := url
	j.ExportURL = &u
	m.Jobs[jobID] = j
	return nil
}

func (m *mockJobStore) GetJob(ctx context.Context, id int) (models.Job, error) {
	j, ok := m.Jobs[id]
	if !ok {
		return models.Job{}, sql.ErrNoRows
	}
	return j, nil
}

func (m *mockJobStore) GetAllJobs(ctx context.Context) ([]models.Job, error) {
	var list []models.Job
	for _, j := range m.Jobs {
		list = append(list, j)
	}
	return list, nil
}

func (m *mockJobStore) DeleteJobs(ctx context.Context, id int) error {
	if _, ok := m.Jobs[id]; !ok {
		return sql.ErrNoRows
	}
	delete(m.Jobs, id)
	return nil
}

func (m *mockJobStore) RecoverStuckJobs(ctx context.Context) (int64, error) {
	var n int64
	for id, j := range m.Jobs {
		if j.Status == pipeline.StatusPending || j.Status == pipeline.StatusRunning {
			j.Status = pipeline.StatusFailed
			m.Jobs[id] = j
			n++
		}
	}
	return n, nil
}

type mockResultStore struct {
	Results map[int][]models.Result
}

func (m *mockResultStore) InsertResults(ctx context.Context, results []models.Result) error {
	for _, r := range results {
		m.Results[r.JobID] = append(m.Results[r.JobID], r)
	}
	return nil
}

func (m *mockResultStore) GetResultsByJob(ctx context.Context, id int) ([]models.Result, error) {
	return m.Results[id], nil
}

type mockErrorStore struct {
	Errors map[int][]models.JobError
}

func (m *mockErrorStore) InsertError(ctx context.Context, err models.JobError) error {
	m.Errors[err.JobID] = append(m.Errors[err.JobID], err)
	return nil
}

func (m *mockErrorStore) GetErrorsByJob(ctx context.Context, id int) ([]models.JobError, error) {
	return m.Errors[id], nil
}

func TestJobService_CreateJob(t *testing.T) {
	js := &mockJobStore{Jobs: make(map[int]models.Job)}
	rs := &mockResultStore{Results: make(map[int][]models.Result)}
	es := &mockErrorStore{Errors: make(map[int][]models.JobError)}
	svc := NewJobService(js, rs, es)

	spec := pipeline.JobSpec{
		Sources: []pipeline.SourceConfig{{Type: "csv", Path: "dummy.csv"}},
	}

	job, err := svc.CreateJob(context.Background(), spec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if job.ID != 1 {
		t.Errorf("expected job ID 1, got %d", job.ID)
	}

	// Wait briefly so background goroutine starts
	// Real pipeline might fail or finish but we just want to see it created.
}

func TestJobService_GetJob(t *testing.T) {
	js := &mockJobStore{Jobs: map[int]models.Job{
		1: {ID: 1, Status: pipeline.StatusRunning},
	}}
	svc := NewJobService(js, nil, nil)

	j, err := svc.GetJob(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if j.ID != 1 {
		t.Errorf("expected job ID 1, got %d", j.ID)
	}
}

func TestJobService_CancelJob(t *testing.T) {
	js := &mockJobStore{Jobs: make(map[int]models.Job)}
	svc := NewJobService(js, nil, nil)

	// Attempting to cancel non-existent/not-running job
	err := svc.CancelJob(context.Background(), 99)
	if err == nil {
		t.Error("expected error when cancelling non-running job")
	}

	js.Jobs[1] = models.Job{ID: 1, Status: pipeline.StatusRunning}
	_, cancel := context.WithCancel(context.Background())
	svc.mu.Lock()
	svc.cancelFuncs[1] = cancel
	svc.mu.Unlock()

	err = svc.CancelJob(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error cancelling job: %v", err)
	}

	if js.Jobs[1].Status != pipeline.StatusCancelled {
		t.Errorf("expected status cancelled, got %s", js.Jobs[1].Status)
	}
}

func TestJobService_DeleteJob(t *testing.T) {
	js := &mockJobStore{Jobs: map[int]models.Job{
		1: {ID: 1, Status: pipeline.StatusRunning},
	}}
	svc := NewJobService(js, nil, nil)

	_, cancel := context.WithCancel(context.Background())
	svc.mu.Lock()
	svc.cancelFuncs[1] = cancel
	svc.mu.Unlock()

	err := svc.DeleteJob(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error deleting job: %v", err)
	}

	if _, ok := js.Jobs[1]; ok {
		t.Errorf("expected job to be deleted from store")
	}

	svc.mu.Lock()
	defer svc.mu.Unlock()
	if _, ok := svc.cancelFuncs[1]; ok {
		t.Errorf("expected cancel func to be removed")
	}
}

func TestJobService_GetAllJobs(t *testing.T) {
	js := &mockJobStore{Jobs: map[int]models.Job{
		1: {ID: 1},
		2: {ID: 2},
	}}
	svc := NewJobService(js, nil, nil)
	jobs, err := svc.GetAllJobs(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 2 {
		t.Errorf("expected 2 jobs, got %d", len(jobs))
	}
}

func TestJobService_GetResults(t *testing.T) {
	rs := &mockResultStore{Results: map[int][]models.Result{
		1: {{JobID: 1, GroupKey: "total", AggregatedValue: 100}},
	}}
	svc := NewJobService(nil, rs, nil)
	res, err := svc.GetResults(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 1 {
		t.Errorf("expected 1 result, got %d", len(res))
	}
}

func TestJobService_GetErrors(t *testing.T) {
	es := &mockErrorStore{Errors: map[int][]models.JobError{
		1: {{JobID: 1, ErrorMessage: "test error"}},
	}}
	svc := NewJobService(nil, nil, es)
	errs, err := svc.GetErrors(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d", len(errs))
	}
}

func TestJobService_RecoverStuckJobs(t *testing.T) {
	js := &mockJobStore{Jobs: map[int]models.Job{
		1: {ID: 1, Status: pipeline.StatusPending},
		2: {ID: 2, Status: pipeline.StatusRunning},
		3: {ID: 3, Status: pipeline.StatusCompleted},
	}}
	svc := NewJobService(js, nil, nil)

	if err := svc.RecoverStuckJobs(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if js.Jobs[1].Status != pipeline.StatusFailed {
		t.Errorf("expected pending job to be marked failed, got %s", js.Jobs[1].Status)
	}
	if js.Jobs[2].Status != pipeline.StatusFailed {
		t.Errorf("expected running job to be marked failed, got %s", js.Jobs[2].Status)
	}
	if js.Jobs[3].Status != pipeline.StatusCompleted {
		t.Errorf("expected completed job to be left alone, got %s", js.Jobs[3].Status)
	}
}
