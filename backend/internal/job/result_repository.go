package job

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/SsagarikaR/pipeline-processing/internal/models"
)

// NewResultStore creates a Postgres-backed ResultStore.
func NewResultStore(db *sql.DB) ResultStore { return &postgresResultStore{db: db} }

// InsertResults saves a job's aggregated results, one row per group.
func (s *postgresResultStore) InsertResults(ctx context.Context, results []models.Result) error {
	for _, r := range results {
		_, err := s.db.ExecContext(ctx,
			`INSERT INTO job_results (job_id, group_key, aggregated_value) VALUES ($1, $2, $3)`,
			r.JobID, r.GroupKey, r.AggregatedValue)
		if err != nil {
			return fmt.Errorf("store: insert result: %w", err)
		}
	}
	return nil
}

// GetResultsByJob returns all the aggregated results for a job.
func (s *postgresResultStore) GetResultsByJob(ctx context.Context, jobID uuid.UUID, limit, offset int) ([]models.Result, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, job_id, group_key, aggregated_value, created_at FROM job_results WHERE job_id = $1 ORDER BY created_at LIMIT $2 OFFSET $3`, jobID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("store: get results: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var out []models.Result
	for rows.Next() {
		var r models.Result
		if err := rows.Scan(&r.ID, &r.JobID, &r.GroupKey, &r.AggregatedValue, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("store: scan result: %w", err)
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
