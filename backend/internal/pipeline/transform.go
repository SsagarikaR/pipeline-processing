package pipeline

import (
	"context"
	"sync"
)

func runTransform(ctx context.Context, jobID int, in <-chan Record, transforms []TransformConfig, workers int, errCh chan<- ProcessError) <-chan Record {
	if workers < 1 {
		workers = 1
	}
	transformedCh := make(chan Record, 100)
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
					out := r
					failed := false
					for _, tc := range transforms {
						fn, ok := GetTransformer(tc.Name)
						if !ok {
							errCh <- ProcessError{JobID: jobID, Stage: "transform", Record: &r, Message: "unknown transform: " + tc.Name}
							failed = true
							break
						}
						var err error
						out, err = fn(out, tc.Params)
						if err != nil {
							errCh <- ProcessError{JobID: jobID, Stage: "transform", Record: &r, Message: err.Error()}
							failed = true
							break
						}
					}
					if failed {
						continue
					}
					select {
					case transformedCh <- out:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(transformedCh)
	}()
	return transformedCh
}