package models

import (
	"encoding/json"
	"time"
)

type Job struct {
	ID                 int            `json:"id"`
	Status             string         `json:"status"`
	Spec               json.RawMessage`json:"spec"`
	TotalRecords       int            `json:"total_records"`
	ProcessedRecords   int            `json:"processed_records"`
	ErrorCount         int            `json:"error_count"`
	CreatedAt          time.Time      `json:"created_at"`
	StartedAt          *time.Time     `json:"started_at"`
	CompletedAt        *time.Time     `json:"completed_at"`
	ExportURL          *string        `json:"export_url"`
}

type Result struct {
	ID                 int       `json:"id"`
	JobID              int       `json:"job_id"`
	GroupKey           string    `json:"group_key"`
	AggregatedValue    float64   `json:"aggregated_value"`
	CreatedAt          time.Time `json:"created_at"`
}

type JobError struct {
	ID           int       `json:"id"`
	JobID        int       `json:"job_id"`
	RecordData   string    `json:"record_data"`
	ErrorMessage string    `json:"error_message"`
	Stage        string    `json:"stage"`
	CreatedAt    time.Time `json:"created_at"`
}