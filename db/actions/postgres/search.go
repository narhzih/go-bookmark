package postgres

import (
	"context"
	"database/sql"
	"github.com/mypipeapp/mypipeapi/db/models"
	"github.com/mypipeapp/mypipeapi/db/repository"
	"github.com/rs/zerolog"
	"time"
)

type searchActions struct {
	Db     *sql.DB
	Logger zerolog.Logger
}

func NewSearchActions(db *sql.DB, logger zerolog.Logger) repository.SearchRepository {
	return searchActions{Db: db, Logger: logger}
}

func (s searchActions) SearchThroughPipes(name string, userId int64) ([]models.Pipe, error) {
	query := `
	SELECT p.id, p.name, p.cover_photo, p.created_at, p.modified_at, p.user_id, COUNT(b.pipe_id) AS total_bookmarks, u.username
	FROM pipes p
		LEFT JOIN bookmarks b ON p.id=b.pipe_id
		LEFT JOIN users u ON p.user_id=u.id
	WHERE 
	    p.user_id=$1
	    AND p.name ILIKE '%' || $2 || '%'
	GROUP BY p.id, u.username
	ORDER BY p.id
    `

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	rows, err := s.Db.QueryContext(ctx, query, userId, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return []models.Pipe{}, ErrNoRecord
		}
		return []models.Pipe{}, err
	}

	var pipes []models.Pipe
	for rows.Next() {
		pipe := models.Pipe{}
		_ = rows.Scan(
			&pipe.ID,
			&pipe.Name,
			&pipe.CoverPhoto,
			&pipe.CreatedAt,
			&pipe.ModifiedAt,
			&pipe.UserID,
			&pipe.Bookmarks,
			&pipe.Creator,
		)
		pipes = append(pipes, pipe)
	}
	return pipes, nil
}

func (s searchActions) SearchThroughTags(name string, userId int64) ([]models.Bookmark, error) {
	query := `
	SELECT
    	bt.bookmark_id, b.user_id, b.pipe_id, b.platform, b.url, b.created_at
	FROM bookmark_tag bt
		INNER JOIN bookmarks b on b.id = bt.bookmark_id
		INNER JOIN tags t on bt.tag_id = t.id
    WHERE
        b.user_id = $1
        AND b.id = bt.bookmark_id
        AND t.name ILIKE '%' || $2 || '%'
    `
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	rows, err := s.Db.QueryContext(ctx, query, userId, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return []models.Bookmark{}, ErrNoRecord
		}
		return []models.Bookmark{}, err
	}
	ba := NewBookmarkActions(s.Db, s.Logger)
	bookmarks := make([]models.Bookmark, 0)
	for rows.Next() {
		bookmark := models.Bookmark{}
		_ = rows.Scan(
			&bookmark.ID,
			&bookmark.UserID,
			&bookmark.PipeID,
			&bookmark.Platform,
			&bookmark.Url,
			&bookmark.CreatedAt,
		)
		bookmark, _ = ba.ParseTags(bookmark)
		bookmarks = append(bookmarks, bookmark)
	}

	return bookmarks, nil
}

func (s searchActions) SearchThroughPlatform(name string, userId int64) ([]models.Bookmark, error) {
	query := `
	SELECT
    	bt.bookmark_id, b.user_id, b.pipe_id, b.platform, b.url, b.created_at
	FROM bookmark_tag bt
		INNER JOIN bookmarks b on b.id = bt.bookmark_id
		INNER JOIN tags t on bt.tag_id = t.id
    WHERE
        b.user_id = $1
        AND b.id = bt.bookmark_id
        AND b.platform ILIKE '%' || $2 || '%'
    `
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	rows, err := s.Db.QueryContext(ctx, query, userId, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return []models.Bookmark{}, ErrNoRecord
		}
		return []models.Bookmark{}, err
	}
	ba := NewBookmarkActions(s.Db, s.Logger)
	bookmarks := make([]models.Bookmark, 0)
	for rows.Next() {
		bookmark := models.Bookmark{}
		_ = rows.Scan(
			&bookmark.ID,
			&bookmark.UserID,
			&bookmark.PipeID,
			&bookmark.Platform,
			&bookmark.Url,
			&bookmark.CreatedAt,
		)
		bookmark, _ = ba.ParseTags(bookmark)
		bookmarks = append(bookmarks, bookmark)
	}

	return bookmarks, nil
}

func (s searchActions) SearchAll(name string, userId int64) ([]interface{}, error) {
	//TODO implement me
	panic("implement me")
}
