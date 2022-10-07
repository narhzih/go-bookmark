package db

import (
	"gitlab.com/trencetech/mypipe-api/db/models"
)

func (db Database) CreateBookmark(bm models.Bookmark) (models.Bookmark, error) {
	var newBm models.Bookmark

	query := "INSERT INTO bookmarks (user_id, pipe_id, platform, url) VALUES($1, $2, $3, $4) RETURNING id, user_id, pipe_id, platform, url"
	err := db.Conn.QueryRow(query, bm.UserID, bm.PipeID, bm.Platform, bm.Url).Scan(
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

func (db Database) GetBookmark(bmID, userID int64) (models.Bookmark, error) {
	var bookmark models.Bookmark

	query := "SELECT id, user_id, pipe_id, platform, url FROM bookmarks WHERE id=$1 AND user_id=$2 LIMIT 1"
	err := db.Conn.QueryRow(query, bmID, userID).Scan(
		&bookmark.ID,
		&bookmark.UserID,
		&bookmark.PipeID,
		&bookmark.Platform,
		&bookmark.Url,
	)
	if err != nil {
		return models.Bookmark{}, err
	}
	return bookmark, nil
}

func (db Database) GetBookmarks(userID, pipeID int64) ([]models.Bookmark, error) {
	var bookmarks []models.Bookmark
	query := "SELECT id, user_id, pipe_id, url, platform FROM bookmarks WHERE user_id=$1 AND pipe_id=$2"
	rows, err := db.Conn.Query(query, userID, pipeID)
	if err != nil {
		return bookmarks, err
	}
	defer rows.Close()

	for rows.Next() {
		var bookmark models.Bookmark
		if err := rows.Scan(&bookmark.ID, &bookmark.UserID, &bookmark.PipeID, &bookmark.Url, &bookmark.Platform); err != nil {
			return bookmarks, err
		}
		bookmarks = append(bookmarks, bookmark)
	}

	if err := rows.Err(); err != nil {
		return bookmarks, err
	}
	return bookmarks, nil
}

func (db Database) GetBookmarksCount(userID int64) (int, error) {
	var bmCount int
	query := "SELECT COUNT(id) FROM bookmarks WHERE user_id=$1"
	err := db.Conn.QueryRow(query, userID).Scan(&bmCount)
	if err != nil {
		return bmCount, err
	}

	return bmCount, nil
}

func (db Database) DeleteBookmark(bmID, userID int64) (bool, error) {
	deleteQuery := "DELETE FROM bookmarks WHERE id=$1 AND user_id=$2"
	_, err := db.Conn.Exec(deleteQuery, bmID, userID)
	if err != nil {
		return false, err
	}
	return true, nil
}
