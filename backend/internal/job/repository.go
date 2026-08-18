package job

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/SsagarikaR/pipeline-processing/internal/models"
)

type JobStore interface {
	CreateJob(ctx context.Context, spec json.RawMessage) (models.Job, error)
	GetJob(ctx context.Context, jobID int) (models.Job, error)
	GetAllJobs(ctx context.Context) ([]models.Job, error)
	DeleteJobs(ctx context.Context, jobID int)(error)
}

type postgresJobStore struct {
	db *sql.DB
}

func NewJobStore(db *sql.DB) JobStore {
	return &postgresJobStore{db: db}
}

func (s *postgresJobStore) CreateJob(ctx context.Context, spec json.RawMessage) (models.Job, error) {
	var job models.Job
	err := s.db.QueryRowContext(ctx, `
	    INSERT INTO jobs (status,spec)
		VALUES ('pending', $1) 
		RETURNING id, status, total_records, processed_records, error_count, created_at, started_at, completed_at
	`,spec).Scan(
		&job.ID,
		&job.Status,
		&job.TotalRecords,
		&job.ProcessedRecords,
		&job.ErrorCount,
		&job.CreatedAt,
		&job.StartedAt,
		&job.CompletedAt,
	)

	if err != nil {
		return job, fmt.Errorf("store: failed to create job: %w", err)
	}
	return job, nil
}

func (s *postgresJobStore) GetJob(ctx context.Context, jobID int) (models.Job, error) {
	var job models.Job
	err := s.db.QueryRowContext(ctx, `
	    SELECT id, status, total_records, processed_records, error_count, created_at, started_at, completed_at
		FROM jobs
		WHERE id = $1
	`,jobID).Scan(
		&job.ID,
		&job.Status,
		&job.TotalRecords,
		&job.ProcessedRecords,
		&job.ErrorCount,
		&job.CreatedAt,
		&job.StartedAt,
		&job.CompletedAt,
	)

	if err != nil {
		return job, fmt.Errorf("store: failed to get job: %w", err)
	}
	return job, nil
}

func (s *postgresJobStore) GetAllJobs(ctx context.Context) ([]models.Job, error) {
	var jobs []models.Job
	rows, err := s.db.QueryContext(ctx, `
	    SELECT id, status, total_records, processed_records, 
	    error_count, created_at, started_at, completed_at
		FROM jobs
	`)

	if err != nil {
		return jobs, fmt.Errorf("store: failed to get all jobs: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var job models.Job
		err := rows.Scan(
			&job.ID,
			&job.Status,
			&job.TotalRecords,
			&job.ProcessedRecords,
			&job.ErrorCount,
			&job.CreatedAt,
			&job.StartedAt,
			&job.CompletedAt,
		)
		if err != nil {
			return jobs, fmt.Errorf("store: failed to scan job: %w", err)
		}
		jobs = append(jobs, job)
	}
	return jobs, rows.Err()
}

func (s *postgresJobStore) DeleteJobs(ctx context.Context,jobID int) (error){
	result, err := s.db.ExecContext(ctx, `DELETE FROM jobs WHERE id = $1`, jobID)
	if err != nil {
		return fmt.Errorf("store delete job %d: %w", jobID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("store check delete result: %w", err)
	}
	if rowsAffected == 0{
		return sql.ErrNoRows
	}
	return  nil
}