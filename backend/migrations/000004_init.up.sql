CREATE INDEX IF NOT EXISTS idx_results_job_id ON job_results(job_id);
CREATE INDEX IF NOT EXISTS idx_job_errors_job_id ON job_errors(job_id);