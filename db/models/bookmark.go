package models

import "time"

type Bookmark struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	PipeID    int64     `json:"pipe_id"`
	Platform  string    `json:"platform"`
	Url       string    `json:"url"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
}
