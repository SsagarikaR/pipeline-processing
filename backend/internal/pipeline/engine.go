package pipeline

import (
	"context"
	
	"github.com/SsagarikaR/pipeline-processing/internal/models"
)

func Run(ctx context.Context, jobID int, spec JobSpec, tracker *Tracker, resultStore func(ctx context.Context, results []models.Result) error) string {
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
			return StatusCancelled
		}
		return StatusCompleted
	case <-ctx.Done():
		return StatusCancelled
	}
}