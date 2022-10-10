package services

import (
	"fmt"
	"gitlab.com/trencetech/mypipe-api/cmd/api/helpers"
	"gitlab.com/trencetech/mypipe-api/db/actions/postgres"
	"gitlab.com/trencetech/mypipe-api/db/models"
)

func (s Services) SharePublicPipe(pipeId, userId int64) (models.SharedPipe, error) {
	var sharedPipeModel models.SharedPipe
	var pipeToBeShared models.Pipe
	var err error

	pipeToBeShared, err = s.Repositories.Pipe.GetPipe(pipeId, userId)
	if err != nil {
		return sharedPipeModel, err
	}
	sharedPipeModel.PipeID = pipeToBeShared.ID
	sharedPipeModel.SharerID = pipeToBeShared.UserID
	sharedPipeModel.Type = "public"
	sharedPipeModel.Code = helpers.RandomToken(15)
	// Parse an empty string to the receiver since it's a public pipe sharer
	sharedPipeModel, err = s.Repositories.PipeShare.CreatePipeShareRecord(sharedPipeModel, "")
	if err != nil {
		return sharedPipeModel, err
	}
	return sharedPipeModel, nil
}

func (s Services) SharePrivatePipe(pipeId, userId int64, shareTo string) (models.SharedPipe, error) {
	var sharedPipeModel models.SharedPipe
	var pipeToBeShared models.Pipe
	var err error

	pipeToBeShared, err = s.Repositories.Pipe.GetPipe(pipeId, userId)
	if err != nil {
		return sharedPipeModel, err
	}
	sharedPipeModel.PipeID = pipeToBeShared.ID
	sharedPipeModel.SharerID = pipeToBeShared.UserID
	sharedPipeModel.Type = "private"
	// Parse an empty string to the receiver since it's a public pipe sharer
	sharedPipeModel, err = s.Repositories.PipeShare.CreatePipeShareRecord(sharedPipeModel, shareTo)
	if err != nil {
		return sharedPipeModel, err
	}
	return sharedPipeModel, nil
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
