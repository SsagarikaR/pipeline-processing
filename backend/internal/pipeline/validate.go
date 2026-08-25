package pipeline

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/uuid"
)

// defaultValidator rejects a record if it has no fields at all.
func defaultValidator(r Record) error {
	if len(r.Data) == 0 {
		return fmt.Errorf("empty record")
	}
	return nil
}

// runValidation is the pipeline's validate stage. A pool of worker
// goroutines checks each record with defaultValidator, valid records
// move on to the next stage, invalid ones are reported as errors (and
// still counted toward progress) instead of being forwarded.
func runValidation(ctx context.Context, jobID uuid.UUID, in <-chan Record, workers int, errCh chan<- ProcessError, progressCh chan<- struct{}) <-chan Record {
	if workers < 1 {
		workers = 1
	}
	validatedCh := make(chan Record, 100)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case r, ok := <-in:
					if !ok {
						return
					}
					if err := defaultValidator(r); err != nil {
						slog.Error("validate stage failed", "job_id", jobID, "error", err)
						select {
						case errCh <- ProcessError{JobID: jobID, Stage: "validate", Record: &r, Message: err.Error()}:
						case <-ctx.Done():
							return
						}

						select {
						case progressCh <- struct{}{}:
						case <-ctx.Done():
							return
						}
						continue
					}
					select {
					case validatedCh <- r:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(validatedCh)
	}()
	return validatedCh
}
