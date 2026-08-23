package job

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/SsagarikaR/pipeline-processing/internal/models"
	"github.com/SsagarikaR/pipeline-processing/internal/pipeline"
)

// NewJobStore creates a Postgres-backed JobStore.
func NewJobStore(db *sql.DB) JobStore {
	return &postgresJobStore{db: db}
}

// CreateJob inserts a new job row with status "pending" and returns the
// row as saved, including its generated ID and timestamps.
func (s *postgresJobStore) CreateJob(ctx context.Context, spec json.RawMessage) (models.Job, error) {
	var job models.Job
	err := s.db.QueryRowContext(ctx, `
	    INSERT INTO jobs (status,spec)
		VALUES ('pending', $1) 
		RETURNING id, status, spec, total_records, processed_records, error_count, created_at, started_at, completed_at, export_url
	`, spec).Scan(
		&job.ID,
		&job.Status,
		&job.Spec,
		&job.TotalRecords,
		&job.ProcessedRecords,
		&job.ErrorCount,
		&job.CreatedAt,
		&job.StartedAt,
		&job.CompletedAt,
		&job.ExportURL,
	)

	if err != nil {
		return job, fmt.Errorf("store: failed to create job: %w", err)
	}
	return job, nil
}

// GetJob fetches one job row by ID.
func (s *postgresJobStore) GetJob(ctx context.Context, jobID uuid.UUID) (models.Job, error) {
	var job models.Job
	err := s.db.QueryRowContext(ctx, `
	    SELECT id, status, spec, total_records, processed_records, error_count, created_at, started_at, completed_at, export_url
		FROM jobs
		WHERE id = $1
	`, jobID).Scan(
		&job.ID,
		&job.Status,
		&job.Spec,
		&job.TotalRecords,
		&job.ProcessedRecords,
		&job.ErrorCount,
		&job.CreatedAt,
		&job.StartedAt,
		&job.CompletedAt,
		&job.ExportURL,
	)

	if err != nil {
		return job, fmt.Errorf("store: failed to get job: %w", err)
	}
	return job, nil
}

// GetAllJobs returns every job row in the table.
func (s *postgresJobStore) GetAllJobs(ctx context.Context) ([]models.Job, error) {
	var jobs []models.Job
	rows, err := s.db.QueryContext(ctx, `
	    SELECT id, status, spec, total_records, processed_records, 
	    error_count, created_at, started_at, completed_at, export_url
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
			&job.ExportURL,
		)
		if err != nil {
			return jobs, fmt.Errorf("store: failed to scan job: %w", err)
		}
		jobs = append(jobs, job)
	}
	return jobs, rows.Err()
}

// DeleteJobs removes a job row by ID. It returns sql.ErrNoRows if no job
// with that ID existed.
func (s *postgresJobStore) DeleteJobs(ctx context.Context, jobID uuid.UUID) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM jobs WHERE id = $1`, jobID)
	if err != nil {
		return fmt.Errorf("store delete job %s: %w", jobID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("store check delete result: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// UpdateStatusAndMetrics writes a job's status and processed/error
// counts, and also stamps started_at, completed_at, or total_records
// depending on which status is being set.
func (s *postgresJobStore) UpdateStatusAndMetrics(ctx context.Context, jobID uuid.UUID, status string, processed int64, errors int64) error {
	query := `UPDATE jobs SET status = $1, processed_records = $2, error_count = $3`
	args := []any{status, processed, errors}

	switch status {
	case pipeline.StatusRunning:
		query += `, started_at = now()`
	case pipeline.StatusCompleted:
		// Once completed, we finally know the total records!
		query += `, total_records = $2, completed_at = now()`
	case pipeline.StatusFailed, pipeline.StatusCancelled:
		query += `, completed_at = now()`
	}
	query += ` WHERE id = $4`
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

// UpdateExportURL saves the URL where a job's exported output landed.
func (s *postgresJobStore) UpdateExportURL(ctx context.Context, jobID uuid.UUID, url string) error {
	result, err := s.db.ExecContext(ctx, `UPDATE jobs SET export_url = $1 WHERE id = $2`, url, jobID)
	if err != nil {
		return fmt.Errorf("store update export url: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("store update export url rows affected: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// RecoverStuckJobs marks every job still sitting in "pending" or
// "running" as failed. Meant to be called once at server startup: if the
// process crashed or was killed mid-run, those jobs have no goroutine
// processing them anymore and would otherwise stay stuck in that status
// forever. Returns how many jobs it recovered.
func (s *postgresJobStore) RecoverStuckJobs(ctx context.Context) (int64, error) {
	result, err := s.db.ExecContext(ctx,
		`UPDATE jobs SET status = $1, completed_at = now() WHERE status IN ($2, $3)`,
		pipeline.StatusFailed, pipeline.StatusPending, pipeline.StatusRunning)
	if err != nil {
		return 0, fmt.Errorf("store: recover stuck jobs: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("store: recover stuck jobs rows affected: %w", err)
	}
	return rows, nil
}
