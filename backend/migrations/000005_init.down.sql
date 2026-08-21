ALTER TABLE jobs DROP COLUMN IF EXISTS export_url;

DROP INDEX IF EXISTS idx_job_errors_job_id_created_at;
CREATE INDEX IF NOT EXISTS idx_job_errors_job_id ON job_errors(job_id);