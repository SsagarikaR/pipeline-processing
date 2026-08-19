package pipeline

import (
	"sync/atomic"

)

type Tracker struct {
	processed  int64
	errCh      chan ProcessError
	progressCh chan struct{}
	onError    func(ProcessError)
}

func NewTracker(onError func(ProcessError)) *Tracker {
	return &Tracker{
		errCh:      make(chan ProcessError, 100),
		progressCh: make(chan struct{}, 1000),
		onError:    onError,
	}
}

func (t *Tracker) Run() {
	go func() {
		for range t.progressCh {
			atomic.AddInt64(&t.processed, 1)
		}
	}()
	go func() {
		for e := range t.errCh {
			t.onError(e)
		}
	}()
}

func (t *Tracker) Processed() int64 { return atomic.LoadInt64(&t.processed) }
func (t *Tracker) Close()           { close(t.progressCh); close(t.errCh) }