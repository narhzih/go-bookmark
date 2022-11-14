package services

import (
	"fmt"
	"github.com/mypipeapp/mypipeapi/cmd/api/helpers"
	"github.com/mypipeapp/mypipeapi/db/actions/postgres"
	"github.com/mypipeapp/mypipeapi/db/models"
)

func (s Services) SharePipePublicly(pipeId, userId int64) (models.SharedPipe, error) {
	var sharedPipeRecord models.SharedPipe
	var pipeToBeShared models.Pipe
	var err error

	pipeToBeShared, err = s.Repositories.Pipe.GetPipe(pipeId, userId)
	if err != nil {
		return sharedPipeRecord, err
	}
	sharedPipeRecord.PipeID = pipeToBeShared.ID
	sharedPipeRecord.SharerID = pipeToBeShared.UserID
	sharedPipeRecord.Type = "public"
	sharedPipeRecord.Code = helpers.RandomToken(15)
	// Parse an empty string to the receiver since it's a public pipe sharer
	sharedPipeRecord, err = s.Repositories.PipeShare.CreatePipeShareRecord(sharedPipeRecord, "")
	if err != nil {
		return sharedPipeRecord, err
	}
	return sharedPipeRecord, nil
}

func (s Services) SharePipePrivately(pipeId, userId int64, shareTo string) (models.SharedPipe, error) {
	var sharedPipeRecord models.SharedPipe
	var pipeToBeShared models.Pipe
	var err error

	pipeToBeShared, err = s.Repositories.Pipe.GetPipe(pipeId, userId)
	if err != nil {
		return sharedPipeRecord, err
	}
	sharedPipeRecord.PipeID = pipeToBeShared.ID
	sharedPipeRecord.SharerID = pipeToBeShared.UserID
	sharedPipeRecord.Type = "private"
	sharedPipeRecord.Code = helpers.RandomToken(15)
	// Parse an empty string to the receiver since it's a public pipe sharer
	sharedPipeRecord, err = s.Repositories.PipeShare.CreatePipeShareRecord(sharedPipeRecord, shareTo)
	if err != nil {
		return sharedPipeRecord, err
	}
	return sharedPipeRecord, nil
}

func (s Services) CanPreviewAndCanAdd(pipe models.Pipe, userId int64) (bool, error) {
	pipeToAdd, err := s.Repositories.PipeShare.GetSharedPipe(pipe.ID)
	if err != nil {
		return false, err
	}
	if pipeToAdd.Type == "private" {
		return false, fmt.Errorf("you can't view pipe because it's a private pipe")
	}
	if pipeToAdd.SharerID == userId {
		return false, postgres.ErrCannotSharePipeToSelf
	}
	return true, nil
}

func (s Services) AddPipeToCollection(pipeToAdd models.SharedPipe, userId int64) (models.SharedPipeReceiver, error) {
	receiverInfo, err := s.Repositories.PipeShare.CreatePipeReceiver(models.SharedPipeReceiver{
		SharedPipeId: pipeToAdd.PipeID,
		ReceiverID:   userId,
	})
	if err != nil {
		return models.SharedPipeReceiver{}, err
	}
	return receiverInfo, nil
}
