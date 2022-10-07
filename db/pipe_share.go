package db

import (
	"database/sql"
	"fmt"
	"github.com/rs/zerolog"
	"gitlab.com/trencetech/mypipe-api/db/models"
	"os"
)

var (
	ErrPipeShareToNotFound   = fmt.Errorf("could not share pipe because user was not found")
	ErrCannotSharePipeToSelf = fmt.Errorf("you cannot share pipe to yourself")
)

func (db Database) CreatePipeShareRecord(pipeShareData models.SharedPipe, receiver string) (models.SharedPipe, error) {
	var query string
	var err error
	if pipeShareData.Type == "public" {
		query = "INSERT INTO shared_pipes (sharer_id, pipe_id, type, code) VALUES ($1, $2, $3, $4) RETURNING id, sharer_id, pipe_id, type, code"
		err = db.Conn.QueryRow(query, pipeShareData.SharerID, pipeShareData.PipeID, pipeShareData.Type, pipeShareData.Code).Scan(
			&pipeShareData.ID,
			&pipeShareData.SharerID,
			&pipeShareData.PipeID,
			&pipeShareData.Type,
			&pipeShareData.Code,
		)
	} else if pipeShareData.Type == "private" {
		pipeShareReceiver, err := db.GetUserByUsername(receiver)
		if err != nil {
			db.Logger.Err(err).Msg("An error occurred while trying to fetch user to share pipe to")

			if err == ErrNoRecord {
				return models.SharedPipe{}, ErrPipeShareToNotFound
			}

			return models.SharedPipe{}, err
		}
		if pipeShareReceiver.ID == pipeShareData.SharerID {
			return models.SharedPipe{}, ErrCannotSharePipeToSelf
		}
		//var sharedTo model.SharedPipeReceiver
		query = `
				INSERT INTO shared_pipes (sharer_id, pipe_id, type) 
				VALUES ($1, $2, $3) 
				RETURNING id, sharer_id, pipe_id, type
		
		`
		err = db.Conn.QueryRow(query, pipeShareData.SharerID, pipeShareData.PipeID, pipeShareData.Type).Scan(
			&pipeShareData.ID,
			&pipeShareData.SharerID,
			&pipeShareData.PipeID,
			&pipeShareData.Type,
		)
		_, err = db.CreatePipeReceiver(models.SharedPipeReceiver{
			SharerId:     pipeShareData.SharerID,
			SharedPipeId: pipeShareData.PipeID,
			ReceiverID:   pipeShareReceiver.ID,
		})
		if err != nil {
			return models.SharedPipe{}, nil
		}
	} else {
		return pipeShareData, fmt.Errorf("invalid pipe type share: %v", pipeShareData.Type)
	}

	if err != nil {
		return pipeShareData, err
	}
	return pipeShareData, nil
}

func (db Database) CreatePipeReceiver(receiver models.SharedPipeReceiver) (models.SharedPipeReceiver, error) {
	query := `
			INSERT INTO shared_pipe_receivers (sharer_id, shared_pipe_id, receiver_id)
			VALUES ($1, $2, $3)
			RETURNING id, sharer_id, shared_pipe_id, receiver_id, created_at, modified_at
	`
	err := db.Conn.QueryRow(query, receiver.SharerId, receiver.SharedPipeId, receiver.ReceiverID).Scan(
		&receiver.ID,
		&receiver.SharerId,
		&receiver.SharedPipeId,
		&receiver.ReceiverID,
		&receiver.CreatedAt,
		&receiver.ModifiedAt,
	)
	if err != nil {
		return receiver, err
	}
	return receiver, nil
}

func (db Database) GetSharedPipe(pipeId int64) (models.SharedPipe, error) {
	var sharedPipe models.SharedPipe
	query := "SELECT id, sharer_id, pipe_id, type, code, created_at FROM shared_pipes WHERE pipe_id=$1 LIMIT 1"
	err := db.Conn.QueryRow(query, pipeId).Scan(
		&sharedPipe.ID,
		&sharedPipe.SharerID,
		&sharedPipe.PipeID,
		&sharedPipe.Type,
		&sharedPipe.Code,
		&sharedPipe.CreatedAt,
	)
	if err != nil {
		if err == ErrNoRecord {
			return models.SharedPipe{}, ErrNoRecord
		}
		return models.SharedPipe{}, err
	}
	return models.SharedPipe{}, nil
}

func (db Database) GetSharedPipeByCode(code string) (models.SharedPipe, error) {
	var sharedPipe models.SharedPipe
	query := "SELECT id,sharer_id, pipe_id, type, code, created_at FROM shared_pipes WHERE code=$1 LIMIT 1"
	err := db.Conn.QueryRow(query, code).Scan(
		&sharedPipe.ID,
		&sharedPipe.SharerID,
		&sharedPipe.PipeID,
		&sharedPipe.Type,
		&sharedPipe.Code,
		&sharedPipe.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.SharedPipe{}, ErrNoRecord
		}
		return models.SharedPipe{}, err
	}
	return sharedPipe, nil
}

func (db Database) GetReceivedPipeRecord(pipeId, userId int64) (models.SharedPipeReceiver, error) {
	logger := zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()

	var sharedPipe models.SharedPipeReceiver
	query := "SELECT id, shared_pipe_id, receiver_id FROM shared_pipe_receivers WHERE shared_pipe_id=$1 AND receiver_id=$2 LIMIT 1"
	err := db.Conn.QueryRow(query, pipeId, userId).Scan(
		&sharedPipe.ID,
		&sharedPipe.SharedPipeId,
		&sharedPipe.ReceiverID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Info().Msg("The logged in user hasn't received this  pipe yet")
			return models.SharedPipeReceiver{}, ErrNoRecord
		}
		logger.Info().Msg("There's a specific error somewhere")
		logger.Err(err)
		return models.SharedPipeReceiver{}, err
	}
	logger.Info().Msg("You have already received this pipe")
	return sharedPipe, nil
}
