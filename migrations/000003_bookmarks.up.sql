CREATE TABLE IF NOT EXISTS bookmarks (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL, 
    pipe_id INT NOT NULL,
    platform VARCHAR(100) NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    FOREIGN KEY (pipe_id) REFERENCES public.pipes  ON DELETE CASCADE
)