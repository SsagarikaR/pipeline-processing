package pipeline

import (
	"context"
	"log/slog"

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
			slog.Error("export stage failed to store results", "job_id", jobID, "error", err)
			errCh <- ProcessError{JobID: jobID, Stage: "export", Message: err.Error()}
		}

		for _, cfg := range exports {
			exporter, ok := GetExporter(cfg.Type)
			if !ok {
				slog.Error("export stage failed: unknown export type", "job_id", jobID, "type", cfg.Type)
				errCh <- ProcessError{JobID: jobID, Stage: "export", Message: "unknown export type: " + cfg.Type}
				continue 
			}
			if err := exporter.Export(ctx, cfg, out.Records, out.Results); err != nil {
				slog.Error("export stage failed", "job_id", jobID, "type", cfg.Type, "error", err)
				errCh <- ProcessError{JobID: jobID, Stage: "export", Message: err.Error()}
			}
		}
	}()

	return doneCh
}