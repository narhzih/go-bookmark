CREATE TABLE IF NOT EXISTS account_verifications (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    token VARCHAR(60) DEFAULT '' UNIQUE,
    expires_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)