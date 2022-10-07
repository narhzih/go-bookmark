package repository

import "gitlab.com/trencetech/mypipe-api/db/models"

type PipeShareRepository interface {
	CreatePipeShareRecord(pipeShareData models.SharedPipe, receiver string) (models.SharedPipe, error)
	CreatePipeReceiver(receiver models.SharedPipeReceiver) (models.SharedPipeReceiver, error)
	GetSharedPipe(pipeId int64) (models.SharedPipe, error)
	GetSharedPipeByCode(code string) (models.SharedPipe, error)
	GetReceivedPipeRecord(pipeId, userId int64) (models.SharedPipeReceiver, error)
}
