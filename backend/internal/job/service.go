package job

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/SsagarikaR/pipeline-processing/internal/models"
)

var ErrInvalidJobType = errors.New("invalid job type")

type JobService struct {
	store JobStore
}

func NewJobService(store JobStore) *JobService {
	return &JobService{store: store}
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
	
	
	return newJob, nil
}

func (s *JobService) GetJob(ctx context.Context, id int) (models.Job, error) {
	return s.store.GetJob(ctx, id)
}

func (s *JobService) GetAllJobs(ctx context.Context) ([]models.Job, error) {
	return s.store.GetAllJobs(ctx)
}

func (s *JobService) DeleteJob(ctx context.Context, id int) error {
	return s.store.DeleteJobs(ctx, id)
}