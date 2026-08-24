package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"
)

type correlationIDKey struct{}

// generateCorrelationID creates a random 16-byte hex string used to tie
// together all the log lines for a single request.
func generateCorrelationID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b)
}

// LoggingMiddleware logs every request as it comes in and again once it
// finishes (with how long it took). It reuses the caller's
// X-Correlation-ID header if one was sent, otherwise generates a new
// one, and echoes it back in the response so callers can trace a
// request end-to-end.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		corrID := r.Header.Get("X-Correlation-ID")
		if corrID == "" {
			corrID = generateCorrelationID()
		}

		ctx := context.WithValue(r.Context(), correlationIDKey{}, corrID)
		r = r.WithContext(ctx)

		slog.Info("incoming request", "method", r.Method, "path", r.URL.Path, "correlation_id", corrID)

		// Ensure response has the header
		w.Header().Set("X-Correlation-ID", corrID)

		next.ServeHTTP(w, r)

		slog.Info("request completed", "method", r.Method, "path", r.URL.Path, "correlation_id", corrID, "duration", time.Since(start))
	})
}
