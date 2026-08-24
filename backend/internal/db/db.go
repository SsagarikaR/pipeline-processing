package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/SsagarikaR/pipeline-processing/internal/config"
)

// New opens a Postgres connection pool, sets sane pool limits, and pings
// the database to make sure it's actually reachable before returning.
func New(ctx context.Context, cfg config.DBConfig) (*sql.DB, error) {
	pool, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("db: open: %w", err)
	}

	pool.SetMaxOpenConns(25)
	pool.SetMaxIdleConns(20)
	pool.SetConnMaxLifetime(30 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.PingContext(pingCtx); err != nil {
		_ = pool.Close()
		return nil, fmt.Errorf("db: ping: %w", err)
	}

	return pool, nil
}
