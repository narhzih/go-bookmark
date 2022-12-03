package repository

import "github.com/mypipeapp/mypipeapi/db/models"

type PipeShareRepository interface {
	CreatePipeShareRecord(pipeShareData models.SharedPipe, receiver string) (models.SharedPipe, error)
	CreatePipeReceiver(receiver models.SharedPipeReceiver) (models.SharedPipeReceiver, error)
	GetSharedPipe(pipeId int64, shareType string) (models.SharedPipe, error)
	GetSharedPipeByCode(code string) (models.SharedPipe, error)
	GetReceivedPipeRecord(pipeId, userId int64) (models.SharedPipeReceiver, error)
	GetReceivedPipeRecordByCode(code string, userId int64) (models.SharedPipeReceiver, error)
	AcceptPrivateShare(receiver models.SharedPipeReceiver) (models.SharedPipeReceiver, error)
}
