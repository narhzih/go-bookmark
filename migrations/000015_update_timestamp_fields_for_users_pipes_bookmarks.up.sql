-- Update users table
ALTER TABLE users
    DROP COLUMN created_at;
ALTER TABLE users
    ADD COLUMN created_at TIMESTAMPTZ DEFAULT now(),
    ADD modified_at TIMESTAMPTZ DEFAULT now();


-- Update pipes table
ALTER TABLE pipes
    DROP COLUMN created_at;
ALTER TABLE pipes
    DROP COLUMN modified_at;
ALTER TABLE pipes
    ADD COLUMN created_at TIMESTAMPTZ DEFAULT now(),
    ADD COLUMN modified_at TIMESTAMPTZ DEFAULT now();

-- Update bookmarks table
ALTER TABLE bookmarks
    DROP COLUMN created_at;
ALTER TABLE bookmarks
    ADD COLUMN created_at TIMESTAMPTZ DEFAULT now();
