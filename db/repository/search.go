package repository

import "github.com/mypipeapp/mypipeapi/db/models"

type SearchRepository interface {
	SearchThroughPipes(name string, userId int64) ([]models.Pipe, error)
	SearchThroughTags(name string, userId int64) ([]models.Bookmark, error)
	SearchThroughPlatform(name string, userId int64) ([]models.Bookmark, error)
	SearchAll(name string, userId int64) ([]interface{}, error)
}
