package job

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/SsagarikaR/pipeline-processing/internal/models"
	"github.com/SsagarikaR/pipeline-processing/internal/pipeline"
)

var ErrInvalidJobType = errors.New("invalid job type")

func NewJobService(store JobStore, resultStore ResultStore, errorStore ErrorStore) *JobService {
	return &JobService{
		store:       store,
		resultStore: resultStore,
		errorStore:  errorStore,
		cancelFuncs: make(map[int]context.CancelFunc),
		trackers:    make(map[int]*pipeline.Tracker),
	}
}

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
func (s *JobService) startJob(jobID int, spec pipeline.JobSpec) {
	runCtx, cancel := context.WithCancel(context.Background())

	tracker := pipeline.NewTracker(func(e pipeline.ProcessError) {
		if err := s.errorStore.InsertError(context.Background(), e.ToJobError()); err != nil {
			fmt.Printf("failed to persist job error: %v\n", err)
		}
	})

	s.mu.Lock()
	s.cancelFuncs[jobID] = cancel
	s.trackers[jobID] = tracker
	s.mu.Unlock()

	if err := s.store.UpdateStatusAndMetrics(context.Background(), jobID, pipeline.StatusRunning, 0, 0); err != nil {
		fmt.Printf("failed to mark job running: %v\n", err)
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
			fmt.Printf("failed to update final status: %v\n", err)
		}
	}()
}

func (s *JobService) GetJob(ctx context.Context, id int) (models.Job, error) {
	return s.store.GetJob(ctx, id)
}

func (s *JobService) GetAllJobs(ctx context.Context) ([]models.Job, error) {
	return s.store.GetAllJobs(ctx)
}

func (s *JobService) DeleteJob(ctx context.Context, id int) error {
	s.mu.Lock()
	if cancel, ok := s.cancelFuncs[id]; ok {
		cancel()
	}
	delete(s.cancelFuncs, id)
	delete(s.trackers, id)
	s.mu.Unlock()

	return s.store.DeleteJobs(ctx, id)
}

func (s *JobService) CancelJob(ctx context.Context, id int) error {
	s.mu.Lock()
	cancel, ok := s.cancelFuncs[id]
	s.mu.Unlock()

	if !ok {
		return errors.New("job is not running")
	}
	cancel()
	return s.store.UpdateStatusAndMetrics(ctx, id, pipeline.StatusCancelled, 0, 0) // Metrics might be overwritten if the pipeline.Run finishes and calls it again, which is fine.
}

func (s *JobService) GetProgress(ctx context.Context, id int) (models.Job, int64, int64, map[string]string, error) {
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

func (s *JobService) GetResults(ctx context.Context, id int) ([]models.Result, error) {
	return s.resultStore.GetResultsByJob(ctx, id)
}

func (s *JobService) GetErrors(ctx context.Context, id int) ([]models.JobError, error) {
	return s.errorStore.GetErrorsByJob(ctx, id)
}
