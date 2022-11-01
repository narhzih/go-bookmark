package models

import "time"

type Pipe struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name,omitempty"`
	UserID     int64     `json:"user_id"`
	CoverPhoto string    `json:"cover_photo"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	Bookmarks  int       `json:"bookmarks"`
	Creator    string    `json:"creator"`
}

type PipeAndResource struct {
	Pipe      Pipe       `json:"pipe"`
	Bookmarks []Bookmark `json:"bookmarks"`
}
