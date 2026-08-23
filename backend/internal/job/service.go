package job

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/SsagarikaR/pipeline-processing/internal/models"
	"github.com/SsagarikaR/pipeline-processing/internal/pipeline"
)

var ErrInvalidJobType = errors.New("invalid job type")

// NewJobService creates a JobService with the given stores and sets up
// the in-memory bookkeeping (cancel functions, progress trackers) it
// needs to manage jobs that are currently running.
func NewJobService(store JobStore, resultStore ResultStore, errorStore ErrorStore) *JobService {
	return &JobService{
		store:       store,
		resultStore: resultStore,
		errorStore:  errorStore,
		cancelFuncs: make(map[uuid.UUID]context.CancelFunc),
		trackers:    make(map[uuid.UUID]*pipeline.Tracker),
	}
}

// CreateJob validates the spec, saves a new job row as "pending", and
// kicks off the pipeline run in the background before returning.
func (s *JobService) CreateJob(ctx context.Context, spec pipeline.JobSpec) (models.Job, error) {
	if len(spec.Sources) == 0 {
		return models.Job{}, ErrInvalidJobType
	}

	rawSpec, err := json.Marshal(spec)
	if err != nil {
		return models.Job{}, fmt.Errorf("marshal spec: %w", err)
	}

	newJob, err := s.store.CreateJob(ctx, rawSpec)
	if err != nil {
		return models.Job{}, err
	}

	s.startJob(newJob.ID, spec)
	return newJob, nil
}

// startJob runs the pipeline for a job in a background goroutine. It
// registers a cancel function and a progress tracker so the job can be
// cancelled or polled for progress while it's running, marks the job as
// "running" up front, and writes the final status once the pipeline
// finishes.
func (s *JobService) startJob(jobID uuid.UUID, spec pipeline.JobSpec) {
	runCtx, cancel := context.WithCancel(context.Background())

	tracker := pipeline.NewTracker(func(e pipeline.ProcessError) {
		if err := s.errorStore.InsertError(context.Background(), e.ToJobError()); err != nil {
			slog.Error("failed to persist job error", "job_id", jobID, "err", err)
		}
	})

	s.mu.Lock()
	s.cancelFuncs[jobID] = cancel
	s.trackers[jobID] = tracker
	s.mu.Unlock()

	if err := s.store.UpdateStatusAndMetrics(context.Background(), jobID, pipeline.StatusRunning, 0, 0); err != nil {
		slog.Error("failed to mark job running", "job_id", jobID, "err", err)
	}

	go func() {
		defer func() {
			s.mu.Lock()
			delete(s.cancelFuncs, jobID)
			s.mu.Unlock()
		}()

		status := pipeline.Run(runCtx, jobID, spec, tracker,
			func(ctx context.Context, results []models.Result) error {
				return s.resultStore.InsertResults(ctx, results)
			},
			func(ctx context.Context, url string) error {
				return s.store.UpdateExportURL(ctx, jobID, url)
			})

		if err := s.store.UpdateStatusAndMetrics(context.Background(), jobID, status, tracker.Processed(), tracker.Errors()); err != nil {
			slog.Error("failed to update final status", "job_id", jobID, "err", err)
		}
	}()
}

// GetJob fetches a single job by ID.
func (s *JobService) GetJob(ctx context.Context, id uuid.UUID) (models.Job, error) {
	return s.store.GetJob(ctx, id)
}

// GetAllJobs lists every job in the system.
func (s *JobService) GetAllJobs(ctx context.Context) ([]models.Job, error) {
	return s.store.GetAllJobs(ctx)
}

// DeleteJob cancels the job if it's still running, forgets its
// in-memory tracking state, and removes it (and its results/errors)
// from the database.
func (s *JobService) DeleteJob(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	if cancel, ok := s.cancelFuncs[id]; ok {
		cancel()
	}
	delete(s.cancelFuncs, id)
	delete(s.trackers, id)
	s.mu.Unlock()

	return s.store.DeleteJobs(ctx, id)
}

// CancelJob stops a running job by calling its cancel function. It
// returns an error if the job isn't currently running (e.g. it already
// finished, or the ID doesn't exist).
func (s *JobService) CancelJob(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	cancel, ok := s.cancelFuncs[id]
	s.mu.Unlock()

	if !ok {
		return errors.New("job is not running")
	}
	cancel()
	return s.store.UpdateStatusAndMetrics(ctx, id, pipeline.StatusCancelled, 0, 0) // Metrics might be overwritten if the pipeline.Run finishes and calls it again, which is fine.
}

// GetProgress returns the job plus its live progress: how many records
// have been processed, how many errors have hit, and per-stage
// latencies. If the job isn't running anymore it falls back to the
// counts stored on the job row itself.
func (s *JobService) GetProgress(ctx context.Context, id uuid.UUID) (models.Job, int64, int64, map[string]string, error) {
	j, err := s.store.GetJob(ctx, id)
	if err != nil {
		return models.Job{}, 0, 0, nil, err
	}

	s.mu.Lock()
	tracker, ok := s.trackers[id]
	s.mu.Unlock()

	var processed, errCount int64
	var stageLatencies map[string]string
	if ok {
		processed = tracker.Processed()
		errCount = tracker.Errors()

		latencies := tracker.StageLatencies()
		stageLatencies = make(map[string]string, len(latencies))
		for k, v := range latencies {
			stageLatencies[k] = v.String()
		}
	} else {
		processed = int64(j.ProcessedRecords)
		errCount = int64(j.ErrorCount)
	}
	return j, processed, errCount, stageLatencies, nil
}

// GetResults returns the aggregated results a job produced.
func (s *JobService) GetResults(ctx context.Context, id uuid.UUID) ([]models.Result, error) {
	return s.resultStore.GetResultsByJob(ctx, id)
}

// GetErrors returns the records that failed processing for a job.
func (s *JobService) GetErrors(ctx context.Context, id uuid.UUID) ([]models.JobError, error) {
	return s.errorStore.GetErrorsByJob(ctx, id)
}

// RecoverStuckJobs marks any job left in "pending" or "running" as
// failed. Call this once at server startup, before accepting traffic:
// those statuses only make sense while a goroutine from a previous
// process is actively driving the job, and after a crash/restart no
// such goroutine exists anymore.
func (s *JobService) RecoverStuckJobs(ctx context.Context) error {
	n, err := s.store.RecoverStuckJobs(ctx)
	if err != nil {
		return fmt.Errorf("recover stuck jobs: %w", err)
	}
	if n > 0 {
		slog.Warn("recovered jobs stuck from a previous run", "count", n)
	}
	return nil
}
