package pipeline

import (
	"context"
	"log/slog"

	"github.com/SsagarikaR/pipeline-processing/internal/models"
)

func Run(ctx context.Context, jobID int, spec JobSpec, tracker *Tracker, resultStore func(ctx context.Context, results []models.Result) error) string {
	slog.Info("pipeline started", "job_id", jobID)
	tracker.Run()
	defer tracker.Close()

	recordsCh := runIngestion(ctx, jobID, spec.Sources, tracker.errCh)
	validatedCh := runValidation(ctx, jobID, recordsCh, spec.Concurrency.ValidateWorkers, tracker.errCh, tracker.progressCh)
	transformedCh := runTransform(ctx, jobID, validatedCh, spec.Transforms, spec.Concurrency.TransformWorkers, tracker.errCh)
	aggCh := runAggregation(ctx, jobID, transformedCh, spec.Aggregations, tracker.progressCh)
	doneCh := runExport(ctx, jobID, aggCh, spec.Exports, resultStore, tracker.errCh)

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