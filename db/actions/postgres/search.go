package postgres

import (
	"database/sql"
	"github.com/mypipeapp/mypipeapi/db/models"
	"github.com/mypipeapp/mypipeapi/db/repository"
	"github.com/rs/zerolog"
)

type searchActions struct {
	Db     *sql.DB
	Logger zerolog.Logger
}

func NewSearchActions(db *sql.DB, logger zerolog.Logger) repository.SearchRepository {
	return searchActions{Db: db, Logger: logger}
}

func (s searchActions) SearchThroughPipes(name string) ([]models.Pipe, error) {
	//TODO implement me
	panic("implement me")
}

func (s searchActions) SearchThroughTags(name string) ([]models.Bookmark, error) {
	//TODO implement me
	panic("implement me")
}

func (s searchActions) SearchAll(name string) ([]interface{}, error) {
	//TODO implement me
	panic("implement me")
}
