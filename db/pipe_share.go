package db

import (
	"database/sql"
	"fmt"
	"github.com/rs/zerolog"
	"gitlab.com/trencetech/mypipe-api/db/model"
	"os"
)

var (
	ErrPipeShareToNotFound   = fmt.Errorf("could not share pipe because user was not found")
	ErrCannotSharePipeToSelf = fmt.Errorf("you cannot share pipe to yourself")
)

func (db Database) CreatePipeShareRecord(pipeShareData model.SharedPipe, receiver string) (model.SharedPipe, error) {
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
				return model.SharedPipe{}, ErrPipeShareToNotFound
			}

			return model.SharedPipe{}, err
		}
		if pipeShareReceiver.ID == pipeShareData.SharerID {
			return model.SharedPipe{}, ErrCannotSharePipeToSelf
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
		_, err = db.CreatePipeReceiver(model.SharedPipeReceiver{
			SharerId:     pipeShareData.SharerID,
			SharedPipeId: pipeShareData.PipeID,
			ReceiverID:   pipeShareReceiver.ID,
		})
		if err != nil {
			return model.SharedPipe{}, nil
		}
	} else {
		return pipeShareData, fmt.Errorf("invalid pipe type share: %v", pipeShareData.Type)
	}

	if err != nil {
		return pipeShareData, err
	}
	return pipeShareData, nil
}

func (db Database) CreatePipeReceiver(receiver model.SharedPipeReceiver) (model.SharedPipeReceiver, error) {
	query := `
			INSERT INTO shared_pipe_receivers (sharer_id, shared_pipe_id, receiver_id, is_accepted)
			VALUES ($1, $2, $3, $4)
			RETURNING id, sharer_id, shared_pipe_id, receiver_id, is_accepted, created_at, modified_at
	`
	err := db.Conn.QueryRow(query, receiver.SharerId, receiver.SharedPipeId, receiver.ReceiverID, receiver.IsAccepted).Scan(
		&receiver.ID,
		&receiver.SharerId,
		&receiver.SharedPipeId,
		&receiver.ReceiverID,
		&receiver.IsAccepted,
		&receiver.CreatedAt,
		&receiver.ModifiedAt,
	)
	if err != nil {
		return receiver, err
	}
	return receiver, nil
}

func (db Database) GetSharedPipe(pipeId int64) (model.SharedPipe, error) {
	var sharedPipe model.SharedPipe
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
			return model.SharedPipe{}, ErrNoRecord
		}
		return model.SharedPipe{}, err
	}
	return model.SharedPipe{}, nil
}

func (db Database) GetSharedPipeByCode(code string) (model.SharedPipe, error) {
	var sharedPipe model.SharedPipe
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
			return model.SharedPipe{}, ErrNoRecord
		}
		return model.SharedPipe{}, err
	}
	return sharedPipe, nil
}

func (db Database) GetReceivedPipeRecord(pipeId, userId int64) (model.SharedPipeReceiver, error) {
	logger := zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()

	var sharedPipe model.SharedPipeReceiver
	query := "SELECT id, shared_pipe_id, receiver_id FROM shared_pipe_receivers WHERE shared_pipe_id=$1 AND receiver_id=$2 LIMIT 1"
	err := db.Conn.QueryRow(query, pipeId, userId).Scan(
		&sharedPipe.ID,
		&sharedPipe.SharedPipeId,
		&sharedPipe.ReceiverID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Info().Msg("The logged in user hasn't received this  pipe yet")
			return model.SharedPipeReceiver{}, ErrNoRecord
		}
		logger.Info().Msg("There's a specific error somewhere")
		logger.Err(err)
		return model.SharedPipeReceiver{}, err
	}
	return sharedPipe, nil
}

func (db Database) AcceptPrivateShare(receiver model.SharedPipeReceiver) (model.SharedPipeReceiver, error) {
	query := `
	UPDATE shared_pipe_receivers 
	SET 
	    is_accepted=true 
	WHERE receiver_id=$1 AND shared_pipe_id=$2
	RETURNING id, sharer_id, shared_pipe_id, receiver_id, is_accepted, created_at, modified_at
	`

	err := db.Conn.QueryRow(query, receiver.ReceiverID, receiver.SharedPipeId).Scan(
		&receiver.ID,
		&receiver.SharerId,
		&receiver.SharedPipeId,
		&receiver.ReceiverID,
		&receiver.IsAccepted,
		&receiver.CreatedAt,
		&receiver.ModifiedAt,
	)
	if err != nil {
		return receiver, err
	}
	return receiver, nil
}
