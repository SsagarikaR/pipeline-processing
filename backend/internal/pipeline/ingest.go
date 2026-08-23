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
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type csvIngester struct{}

// Ingest reads a CSV file (or data URI/URL) row by row, turning each row
// into a Record keyed by the header column names, and streams them out
// on out until the source is exhausted or the context is cancelled.
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

// Ingest reads a JSON source and streams each item out as a Record. It
// accepts either a top-level array of objects, an object that contains
// one array field (that array's items are used), or a single object
// (treated as one record).
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

var SandboxDir = "/data/inputs"

// openSource opens a source path for reading, however it's provided:
// an inline "data:" URI, an "http(s)://" URL (retried a few times with
// backoff), or a local file path restricted to SandboxDir to prevent
// path-traversal reads outside the allowed directory. It returns a
// reader plus a close function the caller must call when done.
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
		var resp *http.Response
		var err error
		// Implement retry logic with exponential backoff for network resilience
		for i := 0; i < 3; i++ {
			resp, err = http.Get(path)
			if err == nil && resp.StatusCode == 200 {
				break
			}
			time.Sleep(time.Duration(1<<i) * time.Second)
		}
		if err != nil {
			return nil, nil, fmt.Errorf("fetch %s: %w", path, err)
		}
		if resp.StatusCode != 200 {
			return nil, nil, fmt.Errorf("fetch %s: bad status %d", path, resp.StatusCode)
		}
		return resp.Body, func() { resp.Body.Close() }, nil
	}
	// Path traversal protection: Clean path and enforce sandbox
	if strings.Contains(path, "..") {
		return nil, nil, fmt.Errorf("open %s: path traversal detected", path)
	}
	cleaned := filepath.Clean(path)
	if SandboxDir != "" && !strings.HasPrefix(cleaned, SandboxDir) {
		return nil, nil, fmt.Errorf("open %s: path outside allowed sandbox %s", path, SandboxDir)
	}
	f, err := os.Open(cleaned)
	if err != nil {
		return nil, nil, fmt.Errorf("open %s: %w", path, err)
	}
	return f, func() { f.Close() }, nil
}

// init registers the built-in ingesters (csv, json, api) so they're
// available as soon as the package is imported. Go runs init functions
// automatically, before main starts.
func init() {
	RegisterIngester("csv", csvIngester{})
	RegisterIngester("json", jsonIngester{})
	RegisterIngester("api", jsonIngester{})
}

// runIngestion is the pipeline's ingest stage: it starts one goroutine
// per configured source, each reading its own records into a shared
// channel, and closes that channel once every source has finished.
func runIngestion(ctx context.Context, jobID uuid.UUID, sources []SourceConfig, errCh chan<- ProcessError) <-chan Record {
	recordsCh := make(chan Record, 100)
	var wg sync.WaitGroup

	for _, src := range sources {
		ingester, ok := GetIngester(src.Type)
		if !ok {
			slog.Error("ingest stage failed: unknown source type", "job_id", jobID, "type", src.Type)
			select {
			case errCh <- ProcessError{JobID: jobID, Stage: "ingest", Message: "unknown source type: " + src.Type}:
			case <-ctx.Done():
			}
			continue
		}
		wg.Add(1)
		go func(cfg SourceConfig, ing Ingester) {
			defer wg.Done()
			if err := ing.Ingest(ctx, cfg, recordsCh); err != nil && err != context.Canceled {
				slog.Error("ingest stage failed", "job_id", jobID, "source", cfg.Path, "error", err)
				select {
				case errCh <- ProcessError{JobID: jobID, Stage: "ingest", Message: err.Error()}:
				case <-ctx.Done():
				}
			}
		}(src, ingester)
	}

	go func() {
		wg.Wait()
		close(recordsCh)
	}()
	return recordsCh
}
