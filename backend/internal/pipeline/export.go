package pipeline

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/SsagarikaR/pipeline-processing/internal/config"
	"github.com/SsagarikaR/pipeline-processing/internal/models"
)

func runExport(ctx context.Context, jobID int, in <-chan aggOutput, exports []ExportConfig, resultStore func(ctx context.Context, results []models.Result) error, storeURL func(ctx context.Context, url string) error, errCh chan<- ProcessError) <-chan struct{} {
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
			url, err := exporter.Export(ctx, cfg, out.Records, out.Results)
			if err != nil {
				slog.Error("export stage failed", "job_id", jobID, "type", cfg.Type, "error", err)
				errCh <- ProcessError{JobID: jobID, Stage: "export", Message: err.Error()}
			}
			if url != "" && storeURL != nil {
				if err := storeURL(ctx, url); err != nil {
					slog.Error("export stage failed to store url", "job_id", jobID, "error", err)
				}
			}
		}
	}()

	return doneCh
}



type s3Exporter struct{}

func (s3Exporter) Export(ctx context.Context, cfg ExportConfig, records []Record, results []models.Result) (string, error) {
	appConfig := config.LoadConfig()
	bucket := appConfig.S3.Bucket
	endpoint := appConfig.S3.Endpoint
	region := appConfig.S3.Region

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
	)
	if err != nil {
		return "", fmt.Errorf("s3 config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	outData := map[string]any{
		"records": records,
		"results": results,
	}
	b, err := json.Marshal(outData)
	if err != nil {
		return "", fmt.Errorf("s3 marshal: %w", err)
	}

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(cfg.Path),
		Body:   bytes.NewReader(b),
	})
	if err != nil {
		return "", fmt.Errorf("s3 put object: %w", err)
	}

	url := fmt.Sprintf("%s/%s/%s", endpoint, bucket, cfg.Path)
	return url, nil
}

func init() {
	RegisterExporter("s3", s3Exporter{})
}