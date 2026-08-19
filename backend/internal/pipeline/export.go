package pipeline

import (
	"context"

	"github.com/SsagarikaR/pipeline-processing/internal/models"
)

func runExport(ctx context.Context, jobID int, in <-chan aggOutput, exports []ExportConfig, resultStore func(ctx context.Context, results []models.Result) error, errCh chan<- ProcessError) <-chan struct{} {
	doneCh := make(chan struct{})

	go func() {
		defer close(doneCh)
		out, ok := <-in
		if !ok {
			return
		}

		if err := resultStore(ctx, out.Results); err != nil {
			errCh <- ProcessError{JobID: jobID, Stage: "export", Message: err.Error()}
		}

		for _, cfg := range exports {
			exporter, ok := GetExporter(cfg.Type)
			if !ok {
				continue 
			}
			if err := exporter.Export(ctx, cfg, out.Records, out.Results); err != nil {
				errCh <- ProcessError{JobID: jobID, Stage: "export", Message: err.Error()}
			}
		}
	}()

	return doneCh
}