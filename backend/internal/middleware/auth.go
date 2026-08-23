package middleware

import (
	"crypto/subtle"
	"net/http"
	"os"
)

// APIKeyMiddleware rejects any request that doesn't carry the correct
// X-API-Key header, except requests to the Swagger UI which stay open
// so the API docs are browsable without a key.
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

		// Constant-time comparison so a mismatch can't be timed to guess
		// the key one byte at a time.
		key := r.Header.Get("X-API-Key")
		if subtle.ConstantTimeCompare([]byte(key), []byte(expectedKey)) != 1 {
			http.Error(w, "401 Unauthorized - Invalid or missing API Key", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
