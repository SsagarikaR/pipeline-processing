package models

import "time"

type Job struct {
	ID                 int         `json:"id"`
	Status             string      `json:"status"`
	TotalRecords       int         `json:"total_records"`
	ProcessedRecords   int         `json:"processed_records"`
	ErrorCount         int         `json:"error_count"`
	CreatedAt          time.Time   `json:"created_at"`
	CompletedAt        *time.Time  `json:"completed_at"` //Pointer because it can be null
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