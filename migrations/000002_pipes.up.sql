CREATE TABLE IF NOT EXISTS pipes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    user_id INT NOT NULL, 
    cover_photo VARCHAR(255) DEFAULT 'NULL',
    created_at TIMESTAMPTZ DEFAULT now(),
    modified_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE (name, user_id)
)