-- populate users table
INSERT INTO users
    (username, profile_name, email, email_verified, cover_photo, twitter_id, device_tokens)
VALUES
    ('user1', 'user1', 'user1@gmail.com', true, 'https://images.unsplash.com/photo-1611608822650-925c227ef4d2?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8Nnx8aGFuZHNvbWUlMjBtYW58ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60', '1234567890', '{"123", "456"}'),
    ('user2', 'user2', 'user2@gmail.com', true, 'https://images.unsplash.com/photo-1611608822650-925c227ef4d2?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8Nnx8aGFuZHNvbWUlMjBtYW58ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60', '', '{}'),
    ('user3', 'user3', 'user3@gmail.com', true, 'https://images.unsplash.com/photo-1611608822650-925c227ef4d2?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8Nnx8aGFuZHNvbWUlMjBtYW58ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60', '', '{}'),
    ('user4', 'user4', 'user4@gmail.com', true, 'https://images.unsplash.com/photo-1611608822650-925c227ef4d2?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8Nnx8aGFuZHNvbWUlMjBtYW58ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60', '', '{}');

-- populate users_auth table
INSERT into user_auth
    (user_id, origin, hashed_password)
VALUES
    (1, 'DEFAULT', '$2a$14$A/CXTnm0.WSb0CoWcH31VeKv.CitRdGTiWHj/06I3cUvwgrj.UwBu'), -- password
    (2, 'DEFAULT', '$2a$14$A/CXTnm0.WSb0CoWcH31VeKv.CitRdGTiWHj/06I3cUvwgrj.UwBu'), -- $passw01
    (3, 'DEFAULT', '$2a$14$A/CXTnm0.WSb0CoWcH31VeKv.CitRdGTiWHj/06I3cUvwgrj.UwBu'), -- $passw01
    (4, 'DEFAULT', '$2a$14$A/CXTnm0.WSb0CoWcH31VeKv.CitRdGTiWHj/06I3cUvwgrj.UwBu'); -- $passw01

-- populate pipes table
INSERT into pipes
    (user_id, name, cover_photo)
VALUES
    (1, 'Youtube Shorts', 'https://images.unsplash.com/photo-1611162616475-46b635cb6868?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8MXx8eW91dHViZSUyMGxvZ298ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60'),
    (1, 'TikTok', 'https://images.unsplash.com/photo-1611605698323-b1e99cfd37ea?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8N3x8eW91dHViZSUyMGxvZ298ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60'),
    (2, 'Youtube Shorts', 'https://images.unsplash.com/photo-1611162616475-46b635cb6868?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8MXx8eW91dHViZSUyMGxvZ298ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60'),
    (2, 'TikTok', 'https://images.unsplash.com/photo-1611605698323-b1e99cfd37ea?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8N3x8eW91dHViZSUyMGxvZ298ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60'),
    (3, 'Youtube Shorts', 'https://images.unsplash.com/photo-1611162616475-46b635cb6868?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8MXx8eW91dHViZSUyMGxvZ298ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60'),
    (3, 'TikTok', 'https://images.unsplash.com/photo-1611605698323-b1e99cfd37ea?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8N3x8eW91dHViZSUyMGxvZ298ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60');

-- populate bookmarks table
INSERT into bookmarks
    (user_id, pipe_id, url, platform)
VALUES
    (1, 1, 'https://youtu.be/Acgk_Jl95es', 'youtube'),
    (1, 2, 'https://www.tiktok.com/@sheebybeauty/video/7159040755863014683?is_from_webapp=1&sender_device=pc', 'tiktok'),
    (2, 3, 'https://youtu.be/Acgk_Jl95es', 'youtube'),
    (2, 4, 'https://www.tiktok.com/@sheebybeauty/video/7159040755863014683?is_from_webapp=1&sender_device=pc', 'tiktok'),
    (3, 5, 'https://youtu.be/Acgk_Jl95es', 'youtube'),
    (3, 6, 'https://www.tiktok.com/@sheebybeauty/video/7159040755863014683?is_from_webapp=1&sender_device=pc', 'tiktok');

-- populate tags table
INSERT INTO tags
    (name)
VALUES
    ('Beautiful Asian Muslim'),
    ('Quick Blows'),
    ('Twerk Videos'),
    ('iOS 16 releases');

-- populate bookmark_tags table
INSERT INTO bookmark_tag
    (bookmark_id, tag_id)
VALUES
    (1, 1),
    (1, 2),
    (2, 2),
    (2, 3),
    (3, 4),
    (3, 1),
    (4, 1),
    (5, 1),
    (6, 2);

-- populate notifications table