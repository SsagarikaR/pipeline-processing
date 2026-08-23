package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAPIKeyMiddleware(t *testing.T) {
	t.Setenv("API_KEY", "test-key")

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := APIKeyMiddleware(next)

	tests := []struct {
		name       string
		path       string
		headerKey  string
		wantStatus int
	}{
		{"correct key is allowed", "/api/v1/pipelines", "test-key", http.StatusOK},
		{"wrong key is rejected", "/api/v1/pipelines", "wrong-key", http.StatusUnauthorized},
		{"missing key is rejected", "/api/v1/pipelines", "", http.StatusUnauthorized},
		{"swagger is allowed without a key", "/swagger/index.html", "", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			if tt.headerKey != "" {
				req.Header.Set("X-API-Key", tt.headerKey)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestAPIKeyMiddleware_DefaultKey(t *testing.T) {
	os.Unsetenv("API_KEY")

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := APIKeyMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pipelines", nil)
	req.Header.Set("X-API-Key", "secret-pipeline-key")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d (default key should work when API_KEY is unset)", rec.Code, http.StatusOK)
	}
}
