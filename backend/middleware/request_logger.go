package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// LoggingMiddleware logs incoming HTTP requests and their duration.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		slog.Info("incoming request", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
		slog.Info("request completed", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))
	})
}
