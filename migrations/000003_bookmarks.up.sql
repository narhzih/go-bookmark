CREATE TABLE IF NOT EXISTS bookmarks (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL, 
    pipe_id INT NOT NULL,
    platform VARCHAR(100) DEFAULT NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    FOREIGN KEY (pipe_id) REFERENCES pipes.id ON DELETE CASCADE
)