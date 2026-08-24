package pipeline

import (
	"context"

	"github.com/SsagarikaR/pipeline-processing/internal/models"
)

type Ingester interface {
	Ingest(ctx context.Context, cfg SourceConfig, out chan<- Record) error
}

type Validator func(r Record) error

type Transformer func(r Record, params map[string]any) (Record, error)

type Exporter interface {
	Export(ctx context.Context, cfg ExportConfig, records []Record, results []models.Result) (string, error)
}

var (
	ingesters    = map[string]Ingester{}
	transformers = map[string]Transformer{}
	exporters    = map[string]Exporter{}
)

// RegisterIngester makes an Ingester available under a source type name
// (e.g. "csv"), so the engine can look it up when a job asks for it.
func RegisterIngester(name string, i Ingester) { ingesters[name] = i }

// RegisterTransformer makes a Transformer available under a name (e.g.
// "uppercase"), so job specs can reference it by that name.
func RegisterTransformer(name string, t Transformer) { transformers[name] = t }

// RegisterExporter makes an Exporter available under an export type name
// (e.g. "s3"), so the engine can look it up when a job asks for it.
func RegisterExporter(name string, e Exporter) { exporters[name] = e }

// GetIngester looks up a registered Ingester by name.
func GetIngester(name string) (Ingester, bool) { i, ok := ingesters[name]; return i, ok }

// GetTransformer looks up a registered Transformer by name.
func GetTransformer(name string) (Transformer, bool) { t, ok := transformers[name]; return t, ok }

// GetExporter looks up a registered Exporter by name.
func GetExporter(name string) (Exporter, bool) { e, ok := exporters[name]; return e, ok }
