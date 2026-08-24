package job

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/SsagarikaR/pipeline-processing/internal/models"
)

// NewErrorStore creates a Postgres-backed ErrorStore.
func NewErrorStore(db *sql.DB) ErrorStore { return &postgresErrorStore{db: db} }

// InsertError saves one record-level failure (e.g. a row that failed
// validation) for a job.
func (s *postgresErrorStore) InsertError(ctx context.Context, e models.JobError) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO job_errors (job_id, record_data, error_message, stage) VALUES ($1, $2, $3, $4)`,
		e.JobID, e.RecordData, e.ErrorMessage, e.Stage)
	if err != nil {
		return fmt.Errorf("store: insert error: %w", err)
	}
	return nil
}

// GetErrorsByJob returns all the errors recorded for a job, oldest first.
func (s *postgresErrorStore) GetErrorsByJob(ctx context.Context, jobID uuid.UUID) ([]models.JobError, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, job_id, record_data, error_message, stage, created_at FROM job_errors WHERE job_id = $1 ORDER BY created_at`, jobID)
	if err != nil {
		return nil, fmt.Errorf("store: get errors: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var out []models.JobError
	for rows.Next() {
		var e models.JobError
		if err := rows.Scan(&e.ID, &e.JobID, &e.RecordData, &e.ErrorMessage, &e.Stage, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("store: scan error: %w", err)
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
