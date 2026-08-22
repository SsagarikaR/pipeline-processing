package middleware

import (
	"net/http"
	"os"
)

// APIKeyMiddleware checks for a valid API Key in the X-API-Key header.
func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow Swagger UI without auth
		if len(r.URL.Path) >= 8 && r.URL.Path[:8] == "/swagger" {
			next.ServeHTTP(w, r)
			return
		}

		// Retrieve expected key from env, fallback to default if not set
		expectedKey := os.Getenv("API_KEY")
		if expectedKey == "" {
			expectedKey = "secret-pipeline-key"
		}

		key := r.Header.Get("X-API-Key")
		if key != expectedKey {
			http.Error(w, "401 Unauthorized - Invalid or missing API Key", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
