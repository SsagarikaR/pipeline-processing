package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/SsagarikaR/pipeline-processing/internal/models"
)

func init() {
	SandboxDir = ""
}

func TestEngineRun_Success(t *testing.T) {
	// Create a temporary JSON file to ingest
	tmpDir := t.TempDir()
	sourcePath := filepath.Join(tmpDir, "input.json")

	data := `[
		{"Height(Inches)": 65.0, "Weight(Pounds)": 150.0},
		{"Height(Inches)": 70.0, "Weight(Pounds)": 180.0}
	]`
	if err := os.WriteFile(sourcePath, []byte(data), 0644); err != nil {
		t.Fatalf("failed to write test input: %v", err)
	}

	spec := JobSpec{
		Sources: []SourceConfig{
			{Type: "json", Path: sourcePath},
		},
		Transforms: []TransformConfig{},
		Aggregations: []AggregationConfig{
			{Field: "Height(Inches)", Op: "avg"},
			{Field: "Weight(Pounds)", Op: "sum"},
		},
		Concurrency: ConcurrencyConfig{
			ValidateWorkers:  2,
			TransformWorkers: 2,
		},
	}

	tracker := NewTracker(func(e ProcessError) {})
	resultStore := func(ctx context.Context, results []models.Result) error {
		// Mock store
		if len(results) != 2 {
			return fmt.Errorf("expected 2 results, got %d", len(results))
		}
		return nil
	}
	urlStore := func(ctx context.Context, url string) error {
		return nil
	}

	status := Run(context.Background(), uuid.New(), spec, tracker, resultStore, urlStore)

	if status != StatusCompleted {
		t.Errorf("expected status %s, got %s", StatusCompleted, status)
	}

	time.Sleep(50 * time.Millisecond)

	if tracker.Processed() != 2 {
		t.Errorf("expected 2 processed records, got %d", tracker.Processed())
	}
	if tracker.Errors() != 0 {
		t.Errorf("expected 0 errors, got %d", tracker.Errors())
	}
}

func TestEngineRun_MalformedSource(t *testing.T) {
	tmpDir := t.TempDir()
	sourcePath := filepath.Join(tmpDir, "input.csv")

	// Missing quotes or malformed CSV to trigger immediate fail
	data := `name,age
"Alice, 25`
	if err := os.WriteFile(sourcePath, []byte(data), 0644); err != nil {
		t.Fatalf("failed to write test input: %v", err)
	}

	spec := JobSpec{
		Sources: []SourceConfig{
			{Type: "csv", Path: sourcePath},
		},
		Concurrency: ConcurrencyConfig{ValidateWorkers: 1, TransformWorkers: 1},
	}

	tracker := NewTracker(func(e ProcessError) {})
	resultStore := func(ctx context.Context, results []models.Result) error {
		return nil
	}
	urlStore := func(ctx context.Context, url string) error {
		return nil
	}

	status := Run(context.Background(), uuid.New(), spec, tracker, resultStore, urlStore)

	if status != StatusFailed {
		t.Errorf("expected status %s, got %s", StatusFailed, status)
	}

	if tracker.Processed() != 0 {
		t.Errorf("expected 0 processed records, got %d", tracker.Processed())
	}
	if tracker.Errors() == 0 {
		t.Errorf("expected >0 errors, got %d", tracker.Errors())
	}
}

func TestEngineRun_API_JSON(t *testing.T) {
	// Mock HTTP Server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Mock the randomuser.me object response structure
		resp := map[string]any{
			"results": []map[string]any{
				{"gender": "female", "nat": "US"},
				{"gender": "male", "nat": "GB"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	spec := JobSpec{
		Sources: []SourceConfig{
			{Type: "json", Path: srv.URL},
		},
		Concurrency: ConcurrencyConfig{ValidateWorkers: 1, TransformWorkers: 1},
	}

	tracker := NewTracker(func(e ProcessError) {})
	resultStore := func(ctx context.Context, results []models.Result) error {
		return nil
	}
	urlStore := func(ctx context.Context, url string) error {
		return nil
	}

	status := Run(context.Background(), uuid.New(), spec, tracker, resultStore, urlStore)

	if status != StatusCompleted {
		t.Errorf("expected status %s, got %s", StatusCompleted, status)
	}

	time.Sleep(50 * time.Millisecond)

	if tracker.Processed() != 2 {
		t.Errorf("expected 2 processed records, got %d", tracker.Processed())
	}
}
