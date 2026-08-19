package pipeline

import (
	"context"
	"fmt"
	"sync"
)

func defaultValidator(r Record) error {
	if len(r.Data) == 0 {
		return fmt.Errorf("empty record")
	}
	return nil
}

func runValidation(ctx context.Context, jobID int, in <-chan Record, workers int, errCh chan<- ProcessError, progressCh chan<- struct{}) <-chan Record {
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
						errCh <- ProcessError{JobID: jobID, Stage: "validate", Record: &r, Message: err.Error()}
						progressCh <- struct{}{}
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