-- Shared pipes can either be public or private
-- The `codes` column will store the code for
-- publicly shared pipes
CREATE TABLE IF NOT EXISTS shared_pipes (
    id SERIAL PRIMARY KEY,
    sharer_id INT NOT NULL,
    pipe_id INT NOT NULL,
    type VARCHAR(50) DEFAULT 'public',
    code VARCHAR(20) DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)