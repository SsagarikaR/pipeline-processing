package pipeline

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
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

	var items []map[string]any
	if err := json.NewDecoder(r).Decode(&items); err != nil {
		return fmt.Errorf("json decode: %w", err)
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
			errCh <- ProcessError{JobID: jobID, Stage: "ingest", Message: "unknown source type: " + src.Type}
			continue
		}
		wg.Add(1)
		go func(cfg SourceConfig, ing Ingester) {
			defer wg.Done()
			if err := ing.Ingest(ctx, cfg, recordsCh); err != nil && err != context.Canceled {
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