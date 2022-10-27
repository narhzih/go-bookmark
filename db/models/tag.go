package models

import "time"

type Tag struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

type BookmarkToTag struct {
	ID         int64  `json:"id"`
	TagId      int64  `json:"tagId"`
	BookmarkId int64  `json:"BookmarkId"`
	TagName    string `json:"tagName"`
}
