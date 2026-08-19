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
	Export(ctx context.Context, cfg ExportConfig, records []Record, results []models.Result) error
}

var (
	ingesters    = map[string]Ingester{}
	transformers = map[string]Transformer{}
	exporters    = map[string]Exporter{}
)

func RegisterIngester(name string, i Ingester)       { ingesters[name] = i }
func RegisterTransformer(name string, t Transformer) { transformers[name] = t }
func RegisterExporter(name string, e Exporter)       { exporters[name] = e }

func GetIngester(name string) (Ingester, bool)       { i, ok := ingesters[name]; return i, ok }
func GetTransformer(name string) (Transformer, bool) { t, ok := transformers[name]; return t, ok }
func GetExporter(name string) (Exporter, bool)       { e, ok := exporters[name]; return e, ok }