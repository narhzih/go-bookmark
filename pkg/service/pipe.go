package service

import (
	"fmt"
	"gitlab.com/trencetech/mypipe-api/db"
)

func (s Service) PipeExists(pipeId, userId int64) (bool, error) {
	pipe, err := s.DB.GetPipe(pipeId, userId)
	if err != nil {
		return false, err
	}
	if pipe.UserID != userId {
		return false, db.ErrNoRecord
	}

	return true, nil
}

func (s Service) UserOwnsPipe(pipeId, userId int64) (bool, error) {
	pipe, err := s.DB.GetPipe(pipeId, userId)
	if err != nil {
		return false, err
	}

	if pipe.UserID != userId {
		return false, fmt.Errorf("pipe does not belong to this user")
	}

	return true, nil
}
