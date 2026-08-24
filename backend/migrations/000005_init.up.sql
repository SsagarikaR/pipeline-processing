ALTER TABLE jobs ADD COLUMN export_url VARCHAR(500);

DROP INDEX IF EXISTS idx_job_errors_job_id;
CREATE INDEX IF NOT EXISTS idx_job_errors_job_id_created_at ON job_errors(job_id, created_at);