package job

import (
	"context"
	"database/sql"
	"encoding/json"
	"sync"

	"github.com/google/uuid"

	"github.com/SsagarikaR/pipeline-processing/internal/models"
	"github.com/SsagarikaR/pipeline-processing/internal/pipeline"
)

type JobStore interface {
	CreateJob(ctx context.Context, spec json.RawMessage) (models.Job, error)
	GetJob(ctx context.Context, jobID uuid.UUID) (models.Job, error)
	GetAllJobs(ctx context.Context, limit, offset int) ([]models.Job, error)
	DeleteJobs(ctx context.Context, jobID uuid.UUID) error
	UpdateStatusAndMetrics(ctx context.Context, jobID uuid.UUID, status string, processed int64, errors int64) error
	UpdateExportURL(ctx context.Context, jobID uuid.UUID, url string) error
	RecoverStuckJobs(ctx context.Context) (int64, error)
}

type postgresJobStore struct {
	db *sql.DB
}

type ResultStore interface {
	InsertResults(ctx context.Context, r []models.Result) error
	GetResultsByJob(ctx context.Context, jobID uuid.UUID, limit, offset int) ([]models.Result, error)
}

type postgresResultStore struct {
	db *sql.DB
}

type ErrorStore interface {
	InsertError(ctx context.Context, e models.JobError) error
	GetErrorsByJob(ctx context.Context, jobID uuid.UUID, limit, offset int) ([]models.JobError, error)
}

type postgresErrorStore struct {
	db *sql.DB
}

type pipelineHandler struct {
	service *JobService
}

type JobService struct {
	store       JobStore
	resultStore ResultStore
	errorStore  ErrorStore

	mu          sync.Mutex
	cancelFuncs map[uuid.UUID]context.CancelFunc
	trackers    map[uuid.UUID]*pipeline.Tracker
}
