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

func NewTracker(onError func(ProcessError)) *Tracker {
	return &Tracker{
		errCh:      make(chan ProcessError, 100),
		progressCh: make(chan struct{}, 1000),
		onError:    onError,
	}
}

func (t *Tracker) RecordStageLatency(stage string, d time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.stageLatencies == nil {
		t.stageLatencies = make(map[string]time.Duration)
	}
	t.stageLatencies[stage] = d
}

func (t *Tracker) StageLatencies() map[string]time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make(map[string]time.Duration, len(t.stageLatencies))
	for k, v := range t.stageLatencies {
		out[k] = v
	}
	return out
}

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

func (t *Tracker) Processed() int64 { return atomic.LoadInt64(&t.processed) }
func (t *Tracker) Errors() int64    { return atomic.LoadInt64(&t.errors) + int64(len(t.errCh)) }
func (t *Tracker) Close()           { close(t.progressCh); close(t.errCh) }
