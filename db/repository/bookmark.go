package repository

import "github.com/mypipeapp/mypipeapi/db/models"

type BookmarkRepository interface {
	CreateBookmark(bm models.Bookmark) (models.Bookmark, error)
	GetBookmark(bmID, userID int64) (models.Bookmark, error)
	GetBookmarks(userID, pipeID int64) ([]models.Bookmark, error)
	ParseTags(bookmark models.Bookmark) (models.Bookmark, error)
	GetBookmarksCount(userID int64) (int, error)
	DeleteBookmark(bmID, userID int64) (bool, error)
}
