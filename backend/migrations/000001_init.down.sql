-- Drop tables in the reverse order of creation to avoid foreign key constraint errors
DROP TABLE IF EXISTS job_errors;
DROP TABLE IF EXISTS job_results;   
DROP TABLE IF EXISTS jobs;
