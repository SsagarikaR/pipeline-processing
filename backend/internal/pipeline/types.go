package pipeline

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/SsagarikaR/pipeline-processing/internal/models"
)

type Record struct {
	Source string
	Data   map[string]any
}

const (
	StatusPending   = "pending"
	StatusRunning   = "running"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
	StatusCancelled = "cancelled"
)

type SourceConfig struct {
	Type string `json:"type"` // "csv" | "json" | "api"
	Path string `json:"path"`
}

type TransformConfig struct {
	Name   string         `json:"name"`
	Params map[string]any `json:"params,omitempty"`
}

type AggregationConfig struct {
	GroupBy string `json:"groupBy,omitempty"`
	Field   string `json:"field"`
	Op      string `json:"op"` // "sum" | "avg" | "count" | "min" | "max"
}

type ExportConfig struct {
	Type string `json:"type"` // "csv" | "json" | "db"
	Path string `json:"path,omitempty"`
}

type ConcurrencyConfig struct {
	ValidateWorkers  int `json:"validateWorkers"`
	TransformWorkers int `json:"transformWorkers"`
}

// JobSpec is what's stored in jobs.spec (json.RawMessage in your store) and
// unmarshaled before a run starts.
type JobSpec struct {
	Sources      []SourceConfig      `json:"sources"`
	Transforms   []TransformConfig   `json:"transforms"`
	Aggregations []AggregationConfig `json:"aggregations"`
	Exports      []ExportConfig      `json:"exports"`
	Concurrency  ConcurrencyConfig   `json:"concurrency"`
}

// ProcessError is the in-flight version of your JobError row - built by
// pipeline stages, then persisted via JobErrorRepository.Insert once the
// error collector picks it up.
type ProcessError struct {
	JobID     uuid.UUID
	Stage     string
	Record    *Record
	Message   string
	CreatedAt time.Time
}

// ToJobError converts an in-flight ProcessError into the models.JobError
// shape that gets saved to the database, serializing the offending
// record's data to JSON if one was attached.
func (e ProcessError) ToJobError() models.JobError {
	data := ""
	if e.Record != nil {
		data = toJSONString(e.Record.Data)
	}
	return models.JobError{
		JobID:        e.JobID,
		RecordData:   data,
		ErrorMessage: e.Message,
		Stage:        e.Stage,
		CreatedAt:    e.CreatedAt,
	}
}

// toJSONString marshals a record's data to a JSON string for storage,
// falling back to a readable error message if it can't be marshaled.
func toJSONString(data map[string]any) string {
	b, err := json.Marshal(data)
	if err != nil {
		return "unable to serialize record: " + err.Error()
	}
	return string(b)
}
