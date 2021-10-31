package model

import "time"

type User struct {
	ID         int64
	Username   string
	Email      string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

type CoverPhoto struct {
	UserID     string
	PhotoUrl   string
	CreatedAt  string
	ModifiedAt string
}

type UserAuth struct {
	User           User
	HashedPassword string
	CreatedAt      time.Time
	ModifiedAt     time.Time
}

type Pipe struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	UserID     int64     `json:"user_id"`
	CoverPhoto string    `json:"cover_photo"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

type Bookmark struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	PipeID    int64     `json:"pipe_id"`
	Platform  string    `json:"platform"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

type PipeAndResource struct {
	Pipe      Pipe       `json:"pipe"`
	Bookmarks []Bookmark `json:"bookmarks"`
}
