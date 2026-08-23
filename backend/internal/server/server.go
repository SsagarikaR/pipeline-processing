package server

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/SsagarikaR/pipeline-processing/internal/config"
	"github.com/SsagarikaR/pipeline-processing/internal/middleware"
)

// New builds the HTTP server: it maps the routes, then wraps them with
// the middleware chain (rate limiting, API key auth, CORS, security
// headers, request logging), outermost middleware first.
func New(cfg *config.Config, pool *sql.DB) *http.Server {
	router := mapRoutes(pool)

	// Create a rate limiter allowing 10 requests per second with a burst of 20
	limiter := middleware.NewRateLimiter(10, 20)

	handler := limiter.Middleware(router)
	handler = middleware.APIKeyMiddleware(handler)
	handler = middleware.CorsMiddleware(cfg.CorsOrigin)(handler)
	handler = middleware.SecurityHeadersMiddleware(handler)
	handler = middleware.LoggingMiddleware(handler)

	return &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

// Shutdown gives the server up to 5 seconds to finish in-flight requests
// before forcing it closed.
func Shutdown(ctx context.Context, srv *http.Server) {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "err", err)
	}
}
