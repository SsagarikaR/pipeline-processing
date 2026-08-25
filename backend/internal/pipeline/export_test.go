package pipeline

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/SsagarikaR/pipeline-processing/internal/models"
)

func TestRunExport(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Register a dummy exporter
	RegisterExporter("dummy", dummyExporter{})

	in := make(chan aggOutput, 1)
	errCh := make(chan ProcessError, 1)
	jobID := uuid.New()

	in <- aggOutput{
		Records: []Record{{Data: map[string]any{"a": 1}}},
		Results: []models.Result{{GroupKey: "k", AggregatedValue: 2}},
	}
	close(in)

	exports := []ExportConfig{{Type: "dummy", Path: "test.json"}}

	resultStoreCalled := false
	resultStore := func(ctx context.Context, results []models.Result) error {
		resultStoreCalled = true
		return nil
	}

	storeURLCalled := false
	storeURL := func(ctx context.Context, url string) error {
		storeURLCalled = true
		if url != "dummy-url" {
			t.Errorf("expected 'dummy-url', got %s", url)
		}
		return nil
	}

	done := runExport(ctx, jobID, in, exports, resultStore, storeURL, errCh)
	<-done

	if !resultStoreCalled {
		t.Error("expected resultStore to be called")
	}
	if !storeURLCalled {
		t.Error("expected storeURL to be called")
	}
}

type dummyExporter struct{}

func (dummyExporter) Export(ctx context.Context, cfg ExportConfig, records []Record, results []models.Result) (string, error) {
	return "dummy-url", nil
}
