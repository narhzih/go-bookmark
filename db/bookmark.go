package db

import "gitlab.com/gowagr/mypipe-api/db/model"

func (db Database) CreateBookmark(bm model.Bookmark) (model.Bookmark, error) {
	var newBm model.Bookmark

	query := "INSERT INTO bookmarks (user_id, pipe_id, platform, url) VALUES($1, $2, $3, $4) RETURNING url"
	err := db.Conn.QueryRow(query, bm.UserID, bm.PipeID, bm.Platform, bm.Url).Scan(
		&newBm.Platform,
		&newBm.Url,
	)

	if err != nil {
		return model.Bookmark{}, err
	}

	return newBm, nil
}

func (db Database) GetBookmark(bmID, userID int64) (model.Bookmark, error) {
	var bookmark model.Bookmark

	return bookmark, nil
}
func (db Database) GetBookmarks(userID int64) ([]model.Bookmark, error) {
	var bookmarks []model.Bookmark
	return bookmarks, nil
}

func (db Database) DeleteBookmark(bmID, userID int64) (bool, error) {
	return true, nil
}
