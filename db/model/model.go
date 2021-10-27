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
	ID         int64
	Name       string
	UserID     int64
	CoverPhoto string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

type Bookmark struct {
	ID        int64
	UserID    int64
	PipeID    int64
	Platform  string
	Url       string
	CreatedAt time.Time
}
