package db

import (
	"fmt"
	"gitlab.com/trencetech/mypipe-api/db/model"
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
			INSERT INTO shared_pipe_receivers (shared_pipe_id, receiver_id)
			VALUES ($1, $2)
			RETURNING id, shared_pipe_id, receiver_id, created_at, modified_at
	`
	err := db.Conn.QueryRow(query, receiver.SharedPipeId, receiver.ReceiverID).Scan(
		&receiver.ID,
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
