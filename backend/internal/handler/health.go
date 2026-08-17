package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type HealthResponse struct {
	Status    string `json:"status"`
	DB        string `json:"db"`
	Timestamp string `json:"timestamp"`
}

func Health(pool *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		overall := "OK"
		dbStatus := "ok"

		pingCtx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := pool.PingContext(pingCtx); err != nil {
			status = http.StatusServiceUnavailable
			overall = "DEGRADED"
			dbStatus = "unreachable"
		}

		response := HealthResponse{
			Status:    overall,
			DB:        dbStatus,
			Timestamp: time.Now().Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			slog.Error("health: failed to write response", "err", err)
		}
	}
}
