package pipeline

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTransformers(t *testing.T) {
	// Lowercase
	fn, ok := GetTransformer("lowercase")
	if !ok {
		t.Fatal("lowercase transformer not found")
	}

	r, err := fn(Record{Data: map[string]any{"text": "HELLO"}}, map[string]any{"field": "text"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if r.Data["text"] != "hello" {
		t.Errorf("expected 'hello', got %v", r.Data["text"])
	}

	// Uppercase
	fn, ok = GetTransformer("uppercase")
	if !ok {
		t.Fatal("uppercase transformer not found")
	}

	r, err = fn(Record{Data: map[string]any{"text": "hello"}}, map[string]any{"field": "text"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if r.Data["text"] != "HELLO" {
		t.Errorf("expected 'HELLO', got %v", r.Data["text"])
	}
}

func TestRunTransform(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	in := make(chan Record, 2)
	errCh := make(chan ProcessError, 2)
	jobID := uuid.New()

	in <- Record{Data: map[string]any{"text": "foo"}}
	close(in)

	transforms := []TransformConfig{
		{Name: "uppercase", Params: map[string]any{"field": "text"}},
	}

	out := runTransform(ctx, jobID, in, transforms, 2, errCh)

	var res []Record
	for r := range out {
		res = append(res, r)
	}

	if len(res) != 1 {
		t.Errorf("expected 1 record, got %d", len(res))
	} else if res[0].Data["text"] != "FOO" {
		t.Errorf("expected 'FOO', got %v", res[0].Data["text"])
	}
}
