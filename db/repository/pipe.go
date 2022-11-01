package repository

import "github.com/mypipeapp/mypipeapi/db/models"

type PipeRepository interface {
	PipeAlreadyExists(pipeName string, userId int64) (bool, error)
	CreatePipe(pipe models.Pipe) (models.Pipe, error)
	GetPipe(pipeId, userId int64) (models.Pipe, error)
	GetPipeByName(pipeName string, userId int64) (models.Pipe, error)
	GetPipeAndResource(pipeId, userId int64) (models.PipeAndResource, error)
	GetPipes(userId int64) ([]models.Pipe, error)
	GetPipesCount(userId int64) (int, error)
	UpdatePipe(userId int64, pipeId int64, updatedBody models.Pipe) (models.Pipe, error)
	DeletePipe(userID, pipeID int64) (bool, error)
}
