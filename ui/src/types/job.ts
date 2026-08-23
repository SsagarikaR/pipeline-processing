// These mirror internal/models/job.go and internal/pipeline/types.go exactly.
// Keep this file in sync if the Go structs change field names.

export type JobStatus = 'pending' | 'running' | 'completed' | 'failed' | 'cancelled';

export interface Job {
  id: number;
  status: JobStatus;
  total_records: number;
  processed_records: number;
  error_count: number;
  created_at: string;
  started_at: string | null;
  completed_at: string | null;
  export_url: string | null;
}

export interface Result {
  id: number;
  job_id: number;
  group_key: string;
  aggregated_value: number;
  created_at: string;
}

export interface JobError {
  id: number;
  job_id: number;
  record_data: string;
  error_message: string;
  stage: string;
  created_at: string;
}

export interface ProgressResponse {
  jobId: number;
  status: JobStatus;
  processed: number;
  errorCount: number;
  percentComplete: number;
  recordsPerSec: number;
  startedAt: string | null;
  completedAt: string | null;
  exportUrl: string | null;
}

// --- JobSpec (the POST /pipelines request body) ---

export type SourceType = 'csv' | 'json';

export interface SourceConfig {
  type: SourceType;
  path: string;
}

export type TransformName = 'uppercase' | 'lowercase';

export interface TransformConfig {
  name: TransformName;
  params: Record<string, string>;
}

export type AggregationOp = 'sum' | 'avg' | 'count' | 'min' | 'max';

export interface AggregationConfig {
  field: string;
  op: AggregationOp;
  groupBy?: string;
}

export type ExportType = 's3' | 'db';

export interface ExportConfig {
  type: ExportType;
  path: string;
}

export interface ConcurrencyConfig {
  validateWorkers: number;
  transformWorkers: number;
}

export interface JobSpec {
  sources: SourceConfig[];
  transforms: TransformConfig[];
  aggregations: AggregationConfig[];
  exports: ExportConfig[];
  concurrency: ConcurrencyConfig;
}

