package service

import (
	"fmt"
	"gitlab.com/trencetech/mypipe-api/db"
	"gitlab.com/trencetech/mypipe-api/db/model"
	"gitlab.com/trencetech/mypipe-api/pkg/helpers"
)

func (s Service) SharePublicPipe(pipeId, userId int64) (model.SharedPipe, error) {
	var sharedPipeModel model.SharedPipe
	var pipeToBeShared model.Pipe
	var err error

	pipeToBeShared, err = s.DB.GetPipe(pipeId, userId)
	if err != nil {
		return sharedPipeModel, err
	}
	sharedPipeModel.PipeID = pipeToBeShared.ID
	sharedPipeModel.SharerID = pipeToBeShared.UserID
	sharedPipeModel.Type = "public"
	sharedPipeModel.Code = helpers.RandomToken(15)
	// Parse an empty string to the receiver since it's a public pipe sharer
	sharedPipeModel, err = s.DB.CreatePipeShareRecord(sharedPipeModel, "")
	if err != nil {
		return sharedPipeModel, err
	}
	return sharedPipeModel, nil
}

func (s Service) SharePrivatePipe(pipeId, userId int64, shareTo string) (model.SharedPipe, error) {
	var sharedPipeModel model.SharedPipe
	var pipeToBeShared model.Pipe
	var err error

	pipeToBeShared, err = s.DB.GetPipe(pipeId, userId)
	if err != nil {
		return sharedPipeModel, err
	}
	sharedPipeModel.PipeID = pipeToBeShared.ID
	sharedPipeModel.SharerID = pipeToBeShared.UserID
	sharedPipeModel.Type = "private"
	// Parse an empty string to the receiver since it's a public pipe sharer
	sharedPipeModel, err = s.DB.CreatePipeShareRecord(sharedPipeModel, shareTo)
	if err != nil {
		return sharedPipeModel, err
	}
	return sharedPipeModel, nil
}

func (s Service) CanPreviewAndCanAdd(pipe model.Pipe, userId int64) (bool, error) {
	pipeToAdd, err := s.DB.GetSharedPipe(pipe.ID)
	if err != nil {
		return false, err
	}
	if pipeToAdd.Type == "private" {
		return false, fmt.Errorf("you can't view pipe because it's a private pipe")
	}
	if pipeToAdd.SharerID == userId {
		return false, db.ErrCannotSharePipeToSelf
	}
	return true, nil
}

func (s Service) AddPipeToCollection(pipeToAdd model.SharedPipe, userId int64) (model.SharedPipeReceiver, error) {
	receiverInfo, err := s.DB.CreatePipeReceiver(model.SharedPipeReceiver{
		SharedPipeId: pipeToAdd.PipeID,
		ReceiverID:   userId,
	})
	if err != nil {
		return model.SharedPipeReceiver{}, err
	}
	return receiverInfo, nil
}
