CREATE TABLE IF NOT EXISTS user_auth
(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    hashed_password VARCHAR(128) NOT NULL
)