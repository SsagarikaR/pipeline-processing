package pipeline

import (
	"context"
	"log/slog"

	"github.com/SsagarikaR/pipeline-processing/internal/models"
)

// Run drives one job through the full pipeline: ingest -> validate ->
// transform -> aggregate -> export. Each stage is its own set of
// goroutines connected by channels, so records stream through rather
// than waiting for earlier stages to fully finish. It blocks until the
// export stage is done (or the context is cancelled) and returns the
// job's final status.
func Run(ctx context.Context, jobID int, spec JobSpec, tracker *Tracker, resultStore func(ctx context.Context, results []models.Result) error, storeURL func(ctx context.Context, url string) error) string {
	slog.Info("pipeline started", "job_id", jobID)
	tracker.Run()
	defer tracker.Close()

	recordsCh := runIngestion(ctx, jobID, spec.Sources, tracker.errCh)
	validatedCh := runValidation(ctx, jobID, recordsCh, spec.Concurrency.ValidateWorkers, tracker.errCh, tracker.progressCh)
	transformedCh := runTransform(ctx, jobID, validatedCh, spec.Transforms, spec.Concurrency.TransformWorkers, tracker.errCh)
	aggCh := runAggregation(ctx, jobID, transformedCh, spec.Aggregations, tracker.progressCh)
	doneCh := runExport(ctx, jobID, aggCh, spec.Exports, resultStore, storeURL, tracker.errCh)

	select {
	case <-doneCh:
		if ctx.Err() != nil {
			slog.Warn("pipeline cancelled", "job_id", jobID)
			return StatusCancelled
		}
		if tracker.Errors() > 0 && tracker.Processed() == 0 {
			slog.Error("pipeline failed: no records processed and errors detected", "job_id", jobID)
			return StatusFailed
		}
		slog.Info("pipeline completed", "job_id", jobID)
		return StatusCompleted
	case <-ctx.Done():
		slog.Warn("pipeline cancelled", "job_id", jobID)
		return StatusCancelled
	}
}
