package postgres

import (
	"database/sql"
	"github.com/mypipeapp/mypipeapi/db/models"
	"github.com/mypipeapp/mypipeapi/db/repository"
	"github.com/rs/zerolog"
)

type bookmarkActions struct {
	Db     *sql.DB
	Logger zerolog.Logger
}

func NewBookmarkActions(db *sql.DB, logger zerolog.Logger) repository.BookmarkRepository {
	return bookmarkActions{
		Db:     db,
		Logger: logger,
	}
}

func (b bookmarkActions) CreateBookmark(bm models.Bookmark) (models.Bookmark, error) {
	var newBm models.Bookmark

	query := "INSERT INTO bookmarks (user_id, pipe_id, platform, url) VALUES($1, $2, $3, $4) RETURNING id, user_id, pipe_id, platform, url"
	err := b.Db.QueryRow(query, bm.UserID, bm.PipeID, bm.Platform, bm.Url).Scan(
		&newBm.ID,
		&newBm.UserID,
		&newBm.PipeID,
		&newBm.Platform,
		&newBm.Url,
	)

	if err != nil {
		return models.Bookmark{}, err
	}

	return newBm, nil
}

func (b bookmarkActions) GetBookmark(bmID, userID int64) (models.Bookmark, error) {
	var bookmark models.Bookmark

	query := "SELECT id, user_id, pipe_id, platform, url, created_at FROM bookmarks WHERE id=$1 AND user_id=$2 LIMIT 1"
	err := b.Db.QueryRow(query, bmID, userID).Scan(
		&bookmark.ID,
		&bookmark.UserID,
		&bookmark.PipeID,
		&bookmark.Platform,
		&bookmark.Url,
		&bookmark.CreatedAt,
	)
	if err != nil {
		return models.Bookmark{}, err
	}
	bookmark, _ = b.ParseTags(bookmark)
	return bookmark, nil
}

func (b bookmarkActions) GetBookmarks(userID, pipeID int64) ([]models.Bookmark, error) {
	var bookmarks []models.Bookmark
	query := "SELECT id, user_id, pipe_id, url, platform, created_at FROM bookmarks WHERE user_id=$1 AND pipe_id=$2"
	rows, err := b.Db.Query(query, userID, pipeID)
	if err != nil {
		return bookmarks, err
	}
	defer rows.Close()

	for rows.Next() {
		var bookmark models.Bookmark
		if err := rows.Scan(
			&bookmark.ID,
			&bookmark.UserID,
			&bookmark.PipeID,
			&bookmark.Url,
			&bookmark.Platform,
			&bookmark.CreatedAt,
		); err != nil {
			return bookmarks, err
		}
		bookmark, _ = b.ParseTags(bookmark)
		bookmarks = append(bookmarks, bookmark)
	}

	if err := rows.Err(); err != nil {
		return bookmarks, err
	}
	return bookmarks, nil
}

func (b bookmarkActions) GetBookmarksCount(userID int64) (int, error) {
	var bmCount int
	query := "SELECT COUNT(id) FROM bookmarks WHERE user_id=$1"
	err := b.Db.QueryRow(query, userID).Scan(&bmCount)
	if err != nil {
		return bmCount, err
	}

	return bmCount, nil
}

func (b bookmarkActions) DeleteBookmark(bmID, userID int64) (bool, error) {
	deleteQuery := "DELETE FROM bookmarks WHERE id=$1 AND user_id=$2"
	_, err := b.Db.Exec(deleteQuery, bmID, userID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (b bookmarkActions) ParseTags(bookmark models.Bookmark) (models.Bookmark, error) {
	query := `
	SELECT bt.id, bt.tag_id, bt.bookmark_id, t.name 
	FROM bookmark_tag bt 
		INNER JOIN tags t on bt.tag_id = t.id
	WHERE bookmark_id=$1
    `

	rows, err := b.Db.Query(query, bookmark.ID)
	if err != nil {
		b.Logger.Err(err).Msg("there was an error parsing tags on bookmark")
		return bookmark, err
	}
	for rows.Next() {
		bookmarkToTag := models.BookmarkToTag{}
		_ = rows.Scan(
			&bookmarkToTag.ID,
			&bookmarkToTag.TagId,
			&bookmarkToTag.BookmarkId,
			&bookmarkToTag.TagName,
		)

		bookmark.Tags = append(bookmark.Tags, bookmarkToTag.TagName)
	}

	return bookmark, nil
}
