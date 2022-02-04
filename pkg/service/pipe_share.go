package service

import (
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
