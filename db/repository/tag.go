package repository

import "github.com/mypipeapp/mypipeapi/db/models"

type TagRepository interface {
	CreateTag(tag models.Tag) (models.Tag, error)
	GetTag(tagId int64) (models.Tag, error)
	GetTagByName(name string) (models.Tag, error)
	AddTagsToBookmark(bmId string, tags []models.Tag) error
}
