ALTER TABLE user_auth
    ADD origin VARCHAR(125) NOT NULL DEFAULT 'DEFAULT';
ALTER TABLE user_auth DROP COLUMN hashed_password;
ALTER TABLE user_auth ADD hashed_password VARCHAR(128) NOT NULL DEFAULT '';