package pipeline

import (
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

type Tracker struct {
	processed      int64
	errors         int64
	errCh          chan ProcessError
	progressCh     chan struct{}
	onError        func(ProcessError)
	mu             sync.Mutex
	stageLatencies map[string]time.Duration
}

// NewTracker creates a Tracker that reports every stage error it sees to
// onError (used to persist errors to the database as they happen).
func NewTracker(onError func(ProcessError)) *Tracker {
	return &Tracker{
		errCh:      make(chan ProcessError, 100),
		progressCh: make(chan struct{}, 1000),
		onError:    onError,
	}
}

// RecordStageLatency saves how long a pipeline stage took, so it can be
// reported back through the progress endpoint.
func (t *Tracker) RecordStageLatency(stage string, d time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.stageLatencies == nil {
		t.stageLatencies = make(map[string]time.Duration)
	}
	t.stageLatencies[stage] = d
}

// StageLatencies returns a copy of the recorded per-stage timings, safe
// to read while the pipeline is still running.
func (t *Tracker) StageLatencies() map[string]time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make(map[string]time.Duration, len(t.stageLatencies))
	for k, v := range t.stageLatencies {
		out[k] = v
	}
	return out
}

// Run starts two background goroutines: one that counts every record
// that finishes processing, and one that counts errors, logs them, and
// forwards each to onError for persistence. Call Close once the
// pipeline is done to stop them.
func (t *Tracker) Run() {
	go func() {
		for range t.progressCh {
			atomic.AddInt64(&t.processed, 1)
		}
	}()
	go func() {
		for e := range t.errCh {
			atomic.AddInt64(&t.errors, 1)
			slog.Error("pipeline stage failed", "job_id", e.JobID, "stage", e.Stage, "error", e.Message)
			t.onError(e)
		}
	}()
}

// Processed returns how many records have finished processing so far.
func (t *Tracker) Processed() int64 { return atomic.LoadInt64(&t.processed) }

// Errors returns how many errors have occurred so far, including any
// still sitting in the error channel waiting to be counted.
func (t *Tracker) Errors() int64 { return atomic.LoadInt64(&t.errors) + int64(len(t.errCh)) }

// Close shuts down the tracker's channels, which stops the goroutines
// started by Run.
func (t *Tracker) Close() { close(t.progressCh); close(t.errCh) }
