CREATE TABLE IF NOT EXISTS shared_pipe_receivers
(
    id SERIAL PRIMARY KEY,
    sharer_id INT NOT NULL,
    shared_pipe_id INT NOT NULL,
    receiver_id INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
)