CREATE TABLE IF NOT EXISTS bookmark_tag (
    id SERIAL PRIMARY KEY,
    bookmark_id INT REFERENCES bookmarks (id) ON DELETE  CASCADE,
    tag_id INT REFERENCES tags (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT now(),
    modified_at TIMESTAMPTZ DEFAULT now()
)