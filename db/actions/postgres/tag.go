package postgres

import (
	"database/sql"
	"github.com/mypipeapp/mypipeapi/db/models"
	"github.com/mypipeapp/mypipeapi/db/repository"
	"github.com/rs/zerolog"
)

type tagActions struct {
	Db     *sql.DB
	Logger zerolog.Logger
}

func NewTagActions(db *sql.DB, logger zerolog.Logger) repository.TagRepository {
	return tagActions{Db: db, Logger: logger}
}

func (t tagActions) CreateTag(tag models.Tag) (models.Tag, error) {
	var createdTag models.Tag
	query := `
	INSERT INTO tags (name)
	VALUES ($1)
	RETURNING id, name, created_at, modified_at
	`

	err := t.Db.QueryRow(query, tag.Name).Scan(
		&tag.Id,
		&tag.Name,
		&tag.CreatedAt,
		&tag.ModifiedAt,
	)
	if err != nil {
		return tag, err
	}
	return createdTag, nil
}

func (t tagActions) GetTag(tagId string) (models.Tag, error) {
	var tag models.Tag
	query := `
	SELECT id, name, created_at, modified_at 
	FROM tags
	WHERE id=$1
    `

	err := t.Db.QueryRow(query, tagId).Scan(
		&tag.Id,
		&tag.Name,
		&tag.CreatedAt,
		&tag.ModifiedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return tag, ErrNoRecord
		}
		return tag, err
	}

	return tag, nil
}
