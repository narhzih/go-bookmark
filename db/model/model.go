package model

import (
	"time"
)

type User struct {
	ID            int64     `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	ProfileName   string    `json:"profile_name"`
	CovertPhoto   string    `json:"cover_photo"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	ModifiedAt    time.Time `json:"modified_at"`
}

type UserAuth struct {
	User           User
	HashedPassword string    `json:"hashed_password"`
	CreatedAt      time.Time `json:"created_at"`
	ModifiedAt     time.Time `json:"modified_at"`
}

type AccountVerification struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	Token      string `json:"token"`
	Used       bool   `json:"used"`
	ExpiresAt  string `json:"expires_at"`
	CreatedAt  string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
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

type Profile struct {
	User      User
	Pipes     int
	Bookmarks int
}
