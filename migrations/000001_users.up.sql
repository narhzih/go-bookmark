CREATE TABLE IF NOT EXISTS users
(
    id SERIAL PRIMARY KEY,
    username VARCHAR(15) NOT NULL UNIQUE,
    profile_name VARCHAR(205) DEFAULT '',
    email VARCHAR(100) DEFAULT '' UNIQUE,
    email_verified BOOLEAN DEFAULT  FALSE,
    device_tokens TEXT ARRAY DEFAULT '{}',
    twitter_id VARCHAR(255) DEFAULT '',
    cover_photo VARCHAR(500) DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT now(),
    modified_at TIMESTAMPTZ DEFAULT now()
);