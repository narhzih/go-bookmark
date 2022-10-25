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

	return createdTag, nil
}
func (t tagActions) GetTag(tagId string) (models.Tag, error) {
	var tag models.Tag

	return tag, nil
}
