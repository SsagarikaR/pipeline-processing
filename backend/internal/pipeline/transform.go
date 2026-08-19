package pipeline

import (
	"context"
	"log/slog"
	"strings"
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
							slog.Error("transform stage failed: unknown transform", "job_id", jobID, "transform", tc.Name)
							errCh <- ProcessError{JobID: jobID, Stage: "transform", Record: &r, Message: "unknown transform: " + tc.Name}
							failed = true
							break
						}
						var err error
						out, err = fn(out, tc.Params)
						if err != nil {
							slog.Error("transform stage failed", "job_id", jobID, "transform", tc.Name, "error", err)
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

func init() {
	RegisterTransformer("lowercase", func(r Record, params map[string]any) (Record, error) {
		if field, ok := params["field"].(string); ok {
			if val, ok := r.Data[field].(string); ok {
				r.Data[field] = strings.ToLower(val)
			}
		}
		return r, nil
	})

	RegisterTransformer("uppercase", func(r Record, params map[string]any) (Record, error) {
		if field, ok := params["field"].(string); ok {
			if val, ok := r.Data[field].(string); ok {
				r.Data[field] = strings.ToUpper(val)
			}
		}
		return r, nil
	})
}