package postgres

import (
	"context"
	"database/sql"
	"github.com/mypipeapp/mypipeapi/db/models"
	"github.com/mypipeapp/mypipeapi/db/repository"
	"github.com/rs/zerolog"
	"time"
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
		&createdTag.Id,
		&createdTag.Name,
		&createdTag.CreatedAt,
		&createdTag.ModifiedAt,
	)
	if err != nil {
		return tag, err
	}
	return createdTag, nil
}

func (t tagActions) GetTag(tagId int64) (models.Tag, error) {
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

func (t tagActions) GetTagByName(name string) (models.Tag, error) {
	var tag models.Tag
	query := `
	SELECT id, name, created_at, modified_at 
	FROM tags
	WHERE name=$1
	LIMIT 1
    `

	err := t.Db.QueryRow(query, name).Scan(
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

func (t tagActions) AddTagsToBookmark(bmId int64, tags []models.Tag) error {
	// This call assumes the bookmark exists already
	for _, tag := range tags {
		eTag, err := t.GetTagByName(tag.Name)
		if err != nil {
			// if there's not existing tag, create one
			if err == ErrNoRecord {
				eTag, err = t.CreateTag(tag)
				if err != nil {
					t.Logger.Err(err).Msg("an error occurred while creating new tag")
				}
			}

			if err != ErrNoRecord {
				t.Logger.Err(err).Msg("some weird type of error")
			}
		}
		query := `
		INSERT INTO bookmark_tag (bookmark_id, tag_id)
		VALUES ($1, $2)
		RETURNING id, bookmark_id, tag_id, created_at, modified_at
		`

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		_, err = t.Db.ExecContext(ctx, query, bmId, eTag.Id)
		if err != nil {
			t.Logger.Err(err).Msg("an error occurred while creating new tag")
		}
	}
	return nil
}
