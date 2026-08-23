package job

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/SsagarikaR/pipeline-processing/internal/models"
)

func NewResultStore(db *sql.DB) ResultStore { return &postgresResultStore{db: db} }

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

func (s *postgresResultStore) GetResultsByJob(ctx context.Context, jobID int) ([]models.Result, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, job_id, group_key, aggregated_value, created_at FROM job_results WHERE job_id = $1`, jobID)
	if err != nil {
		return nil, fmt.Errorf("store: get results: %w", err)
	}
	defer rows.Close()

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
