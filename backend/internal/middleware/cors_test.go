package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCorsMiddleware(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := CorsMiddleware("http://localhost:5173")(next)

	t.Run("echoes the configured origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/pipelines", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
			t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, "http://localhost:5173")
		}
	})

	t.Run("allows every HTTP method the API actually uses", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/api/v1/pipelines/1/cancel", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		allowed := rec.Header().Get("Access-Control-Allow-Methods")
			for _, method := range []string{"GET", "POST", "PATCH", "DELETE"} {
			if !strings.Contains(allowed, method) {
				t.Errorf("Access-Control-Allow-Methods = %q, missing %q", allowed, method)
			}
		}
	})

	t.Run("short-circuits OPTIONS preflight with 200 and no next handler call", func(t *testing.T) {
		called := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true })
		h := CorsMiddleware("http://localhost:5173")(next)

		req := httptest.NewRequest(http.MethodOptions, "/api/v1/pipelines", nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if called {
			t.Error("expected the next handler not to be called for an OPTIONS preflight")
		}
	})
}
