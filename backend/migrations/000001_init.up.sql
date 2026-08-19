-- Create a table to store pipeline job metadata
CREATE TABLE jobs (
    id SERIAL PRIMARY KEY,              
    status VARCHAR(50) NOT NULL,        
    total_records INT DEFAULT 0,       
    processed_records INT DEFAULT 0,    
    error_count INT DEFAULT 0,          
    created_at TIMESTAMP DEFAULT NOW(), 
    completed_at TIMESTAMP              
);

-- Create a table to store individual processed results
CREATE TABLE job_results (
    id SERIAL PRIMARY KEY,             
    job_id INT NOT NULL,              
    group_key VARCHAR(255) NOT NULL,    
    aggregated_value FLOAT NOT NULL,    
    created_at TIMESTAMP DEFAULT NOW(), 
    CONSTRAINT fk_job FOREIGN KEY(job_id) REFERENCES jobs(id) ON DELETE CASCADE 
);

-- Create a table to store errors and failed records for a job
CREATE TABLE job_errors (
    id SERIAL PRIMARY KEY,
    job_id INT NOT NULL,
    record_data TEXT,                 
    error_message TEXT NOT NULL,        
    stage VARCHAR(50),                  
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_job FOREIGN KEY(job_id) REFERENCES jobs(id) ON DELETE CASCADE
);
