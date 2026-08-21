package pipeline

import (
	"context"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
)

type csvIngester struct{}

func (csvIngester) Ingest(ctx context.Context, cfg SourceConfig, out chan<- Record) error {
	r, closeFn, err := openSource(cfg.Path)
	if err != nil {
		return err
	}
	defer closeFn()

	reader := csv.NewReader(r)
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("csv header: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		row, err := reader.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("csv row: %w", err)
		}
		data := make(map[string]any, len(header))
		for i, col := range header {
			if i < len(row) {
				data[col] = row[i]
			}
		}
		select {
		case out <- Record{Source: cfg.Path, Data: data}:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

type jsonIngester struct{}

func (jsonIngester) Ingest(ctx context.Context, cfg SourceConfig, out chan<- Record) error {
	r, closeFn, err := openSource(cfg.Path)
	if err != nil {
		return err
	}
	defer closeFn()

	var raw any
	if err := json.NewDecoder(r).Decode(&raw); err != nil {
		return fmt.Errorf("json decode: %w", err)
	}

	var items []map[string]any
	switch v := raw.(type) {
	case []any:
		for _, item := range v {
			if m, ok := item.(map[string]any); ok {
				items = append(items, m)
			}
		}
	case map[string]any:
			for _, val := range v {
			if arr, ok := val.([]any); ok {
				for _, item := range arr {
					if m, ok := item.(map[string]any); ok {
						items = append(items, m)
					}
				}
				break 
			}
		}
		if len(items) == 0 {
			items = append(items, v)
		}
	default:
		return fmt.Errorf("json decode: expected array or object, got %T", raw)
	}

	for _, item := range items {
		select {
		case out <- Record{Source: cfg.Path, Data: item}:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

func openSource(path string) (io.Reader, func(), error) {
	if strings.HasPrefix(path, "data:") {
		parts := strings.SplitN(path, ",", 2)
		if len(parts) != 2 {
			return nil, nil, fmt.Errorf("invalid data url format")
		}
		
		var r io.Reader
		if strings.Contains(parts[0], ";base64") {
			r = base64.NewDecoder(base64.StdEncoding, strings.NewReader(parts[1]))
		} else {
			r = strings.NewReader(parts[1])
		}
		return r, func() {}, nil
	}
	if strings.HasPrefix(path, "http") {
		resp, err := http.Get(path)
		if err != nil {
			return nil, nil, fmt.Errorf("fetch %s: %w", path, err)
		}
		return resp.Body, func() { resp.Body.Close() }, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("open %s: %w", path, err)
	}
	return f, func() { f.Close() }, nil
}


func init() { //init is special function in go it automatically runs when the package loaded before anything else runs
	RegisterIngester("csv", csvIngester{})
	RegisterIngester("json", jsonIngester{})
	RegisterIngester("api", jsonIngester{})
}

func runIngestion(ctx context.Context, jobID int, sources []SourceConfig, errCh chan<- ProcessError) <-chan Record {
	recordsCh := make(chan Record, 100)
	var wg sync.WaitGroup

	for _, src := range sources {
		ingester, ok := GetIngester(src.Type)
		if !ok {
			slog.Error("ingest stage failed: unknown source type", "job_id", jobID, "type", src.Type)
			errCh <- ProcessError{JobID: jobID, Stage: "ingest", Message: "unknown source type: " + src.Type}
			continue
		}
		wg.Add(1)
		go func(cfg SourceConfig, ing Ingester) {
			defer wg.Done()
			if err := ing.Ingest(ctx, cfg, recordsCh); err != nil && err != context.Canceled {
				slog.Error("ingest stage failed", "job_id", jobID, "source", cfg.Path, "error", err)
				errCh <- ProcessError{JobID: jobID, Stage: "ingest", Message: err.Error()}
			}
		}(src, ingester)
	}

	go func() {
		wg.Wait()
		close(recordsCh)
	}()
	return recordsCh
}