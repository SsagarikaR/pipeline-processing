package job

import (
	"context"
	"database/sql"
	"encoding/json"
	"sync"

	"github.com/SsagarikaR/pipeline-processing/internal/models"
	"github.com/SsagarikaR/pipeline-processing/internal/pipeline"
)

type JobStore interface {
	CreateJob(ctx context.Context, spec json.RawMessage) (models.Job, error)
	GetJob(ctx context.Context, jobID int) (models.Job, error)
	GetAllJobs(ctx context.Context) ([]models.Job, error)
	DeleteJobs(ctx context.Context, jobID int) error
	UpdateStatusAndMetrics(ctx context.Context, jobID int, status string, processed int64, errors int64) error
	UpdateExportURL(ctx context.Context, jobID int, url string) error
	RecoverStuckJobs(ctx context.Context) (int64, error)
}

type postgresJobStore struct {
	db *sql.DB
}

type ResultStore interface {
	InsertResults(ctx context.Context, r []models.Result) error
	GetResultsByJob(ctx context.Context, jobID int) ([]models.Result, error)
}

type postgresResultStore struct {
	db *sql.DB
}

type ErrorStore interface {
	InsertError(ctx context.Context, e models.JobError) error
	GetErrorsByJob(ctx context.Context, jobID int) ([]models.JobError, error)
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
	cancelFuncs map[int]context.CancelFunc
	trackers    map[int]*pipeline.Tracker
}
