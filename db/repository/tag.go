package repository

import "github.com/mypipeapp/mypipeapi/db/models"

type TagRepository interface {
	CreateTag(tag models.Tag) (models.Tag, error)
	GetTag(tagId string) (models.Tag, error)
}
