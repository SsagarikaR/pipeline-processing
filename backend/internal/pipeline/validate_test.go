package pipeline

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestDefaultValidator(t *testing.T) {
	err := defaultValidator(Record{Data: map[string]any{}})
	if err == nil {
		t.Error("expected error for empty record")
	}

	err = defaultValidator(Record{Data: map[string]any{"key": "value"}})
	if err != nil {
		t.Error("expected no error for valid record")
	}
}

func TestRunValidation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	in := make(chan Record, 2)
	errCh := make(chan ProcessError, 2)
	progressCh := make(chan struct{}, 2)
	jobID := uuid.New()

	in <- Record{Data: map[string]any{"valid": true}}
	in <- Record{Data: map[string]any{}} // invalid
	close(in)

	out := runValidation(ctx, jobID, in, 2, errCh, progressCh)

	validCount := 0
	for range out {
		validCount++
	}

	if validCount != 1 {
		t.Errorf("expected 1 valid record, got %d", validCount)
	}

	select {
	case err := <-errCh:
		if err.Stage != "validate" {
			t.Errorf("expected validate stage error, got %s", err.Stage)
		}
	default:
		t.Error("expected an error on errCh")
	}
}
