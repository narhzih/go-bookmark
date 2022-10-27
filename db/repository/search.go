package repository

import "github.com/mypipeapp/mypipeapi/db/models"

type SearchRepository interface {
	SearchThroughPipes(name string) ([]models.Pipe, error)
	SearchThroughTags(name string) ([]models.Bookmark, error)
	SearchAll(name string) ([]interface{}, error)
}
