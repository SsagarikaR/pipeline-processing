package job

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/SsagarikaR/pipeline-processing/internal/models"
	"github.com/SsagarikaR/pipeline-processing/internal/pipeline"
	
)

type JobStore interface {
	CreateJob(ctx context.Context, spec json.RawMessage) (models.Job, error)
	GetJob(ctx context.Context, jobID int) (models.Job, error)
	GetAllJobs(ctx context.Context) ([]models.Job, error)
	DeleteJobs(ctx context.Context, jobID int)(error)
	UpdateStatus(ctx context.Context, jobID int,status string)(error)
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
		RETURNING id, status, spec, total_records, processed_records, error_count, created_at, started_at, completed_at
	`,spec).Scan(
		&job.ID,
		&job.Status,
		&job.Spec,
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
	    SELECT id, status, spec, total_records, processed_records, error_count, created_at, started_at, completed_at
		FROM jobs
		WHERE id = $1
	`,jobID).Scan(
		&job.ID,
		&job.Status,
		&job.Spec,
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
	    SELECT id, status, spec, total_records, processed_records, 
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
			&job.Spec,
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

func (s *postgresJobStore) UpdateStatus(ctx context.Context, jobID int, status string) error {
	query := `UPDATE jobs SET status = $1`
	args := []any{status}

	// Stamp started_at/completed_at as the job transitions, so GET /jobs/:id
	// and /progress can report real timestamps without extra plumbing.
	switch status {
	case pipeline.StatusRunning:
		query += `, started_at = now()`
	case pipeline.StatusCompleted, pipeline.StatusFailed, pipeline.StatusCancelled:
		query += `, completed_at = now()`
	}
	query += ` WHERE id = $2`
	args = append(args, jobID)

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("store: update status: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("store: update status rows affected: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}