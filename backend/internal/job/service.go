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

func (s *JobService) CreateJob(ctx context.Context, jobType string, data []int) (models.Job, error) {
	if jobType == "" {
		return models.Job{}, ErrInvalidJobType
	}
	if len(data) == 0 {
		return models.Job{}, errors.New("data cannot be empty")
	}

	spec, err := json.Marshal(struct {
		Type string `json:"type"`
		Data []int  `json:"data"`
	}{Type: jobType, Data: data})
	if err != nil {
		return models.Job{}, err
	}

	return s.store.CreateJob(ctx, spec)
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