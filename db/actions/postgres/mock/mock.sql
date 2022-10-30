-- populate users table
INSERT INTO users
    (username, profile_name, email, email_verified, cover_photo)
VALUES
    ('user1', 'user1', 'user1@gmail.com', true, 'https://images.unsplash.com/photo-1611608822650-925c227ef4d2?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8Nnx8aGFuZHNvbWUlMjBtYW58ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60'),
    ('user2', 'user2', 'user2@gmail.com', true, 'https://images.unsplash.com/photo-1611608822650-925c227ef4d2?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8Nnx8aGFuZHNvbWUlMjBtYW58ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60'),
    ('user3', 'user3', 'user3@gmail.com', true, 'https://images.unsplash.com/photo-1611608822650-925c227ef4d2?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8Nnx8aGFuZHNvbWUlMjBtYW58ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60'),
    ('user4', 'user4', 'user4@gmail.com', true, 'https://images.unsplash.com/photo-1611608822650-925c227ef4d2?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8Nnx8aGFuZHNvbWUlMjBtYW58ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60');

-- populate users_auth table
INSERT into user_auth
    (user_id, origin, hashed_password)
VALUES
    (1, 'DEFAULT', '$2y$15$salteadoususuueyryy28u48viMdUKIwgSc.ETLYvODrrv3MFczPq'), -- $passw01
    (2, 'DEFAULT', '$2y$15$salteadoususuueyryy28u48viMdUKIwgSc.ETLYvODrrv3MFczPq'), -- $passw01
    (3, 'DEFAULT', '$2y$15$salteadoususuueyryy28u48viMdUKIwgSc.ETLYvODrrv3MFczPq'), -- $passw01
    (4, 'DEFAULT', '$2y$15$salteadoususuueyryy28u48viMdUKIwgSc.ETLYvODrrv3MFczPq'); -- $passw01

-- populate pipes table

-- populate bookmarks table

-- populate tags table

-- populate bookmark_tags table

-- populate notifications table