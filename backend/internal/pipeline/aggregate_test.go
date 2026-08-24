package pipeline

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestRunAggregation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	in := make(chan Record, 3)
	progressCh := make(chan struct{}, 3)
	jobID := uuid.New()

	in <- Record{Data: map[string]any{"dept": "sales", "amount": 100}}
	in <- Record{Data: map[string]any{"dept": "sales", "amount": 200}}
	in <- Record{Data: map[string]any{"dept": "hr", "amount": 50}}
	close(in)

	configs := []AggregationConfig{
		{Op: "sum", Field: "amount", GroupBy: "dept"},
		{Op: "avg", Field: "amount", GroupBy: "dept"},
		{Op: "count", Field: "amount"}, // total count
	}

	outCh := runAggregation(ctx, jobID, in, configs, progressCh)
	out := <-outCh

	if len(out.Records) != 3 {
		t.Errorf("expected 3 records, got %d", len(out.Records))
	}

	if len(out.Results) != 5 { // sum:amount:sales, sum:amount:hr, avg:amount:sales, avg:amount:hr, count:amount
		t.Errorf("expected 5 results, got %d", len(out.Results))
	}
}

func TestToFloat(t *testing.T) {
	v, ok := toFloat(100)
	if !ok || v != 100.0 {
		t.Errorf("failed int to float")
	}

	v, ok = toFloat("100.5")
	if !ok || v != 100.5 {
		t.Errorf("failed string to float")
	}
}
