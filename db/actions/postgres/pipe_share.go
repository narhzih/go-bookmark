package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mypipeapp/mypipeapi/db/models"
	"github.com/mypipeapp/mypipeapi/db/repository"
	"github.com/rs/zerolog"
	"os"
	"time"
)

var (
	ErrPipeShareToNotFound   = fmt.Errorf("could not share pipe because user was not found")
	ErrCannotSharePipeToSelf = fmt.Errorf("you cannot share pipe to yourself")
)

type pipeShareActions struct {
	Db     *sql.DB
	Logger zerolog.Logger
}

func NewPipeShareActions(db *sql.DB, logger zerolog.Logger) repository.PipeShareRepository {
	return pipeShareActions{
		Db:     db,
		Logger: logger,
	}
}

// CreatePipeShareRecord creates a pipe share record for a user
func (p pipeShareActions) CreatePipeShareRecord(pipeShareData models.SharedPipe, receiver string) (models.SharedPipe, error) {
	var query string
	var err error
	uActions := NewUserActions(p.Db, p.Logger)
	switch pipeShareData.Type {
	case models.PipeShareTypePublic:
		query = `
		INSERT INTO shared_pipes 
		    (sharer_id, pipe_id, type, code) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, sharer_id, pipe_id, type, code, created_at, modified_at
		`

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		err = p.Db.QueryRowContext(ctx, query, pipeShareData.SharerID, pipeShareData.PipeID, pipeShareData.Type, pipeShareData.Code).Scan(
			&pipeShareData.ID,
			&pipeShareData.SharerID,
			&pipeShareData.PipeID,
			&pipeShareData.Type,
			&pipeShareData.Code,
			&pipeShareData.CreatedAt,
			&pipeShareData.ModifiedAt,
		)
	case models.PipeShareTypePrivate:
		pipeShareReceiver, err := uActions.GetUserByUsername(receiver)
		if err != nil {
			p.Logger.Err(err).Msg("An error occurred while trying to fetch user to share pipe to")

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
				INSERT INTO shared_pipes (sharer_id, pipe_id, type, code) 
				VALUES ($1, $2, $3, $4) 
				RETURNING id, sharer_id, pipe_id, type, code, created_at, modified_at
		
		`

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		err = p.Db.QueryRowContext(ctx, query, pipeShareData.SharerID, pipeShareData.PipeID, pipeShareData.Type, pipeShareData.Code).Scan(
			&pipeShareData.ID,
			&pipeShareData.SharerID,
			&pipeShareData.PipeID,
			&pipeShareData.Type,
			&pipeShareData.Code,
			&pipeShareData.CreatedAt,
			&pipeShareData.ModifiedAt,
		)

		if err != nil {
			return models.SharedPipe{}, err
		}
		_, err = p.CreatePipeReceiver(models.SharedPipeReceiver{
			SharerId:     pipeShareData.SharerID,
			SharedPipeId: pipeShareData.PipeID,
			ReceiverID:   pipeShareReceiver.ID,
		})
		if err != nil {
			return models.SharedPipe{}, err
		}
	default:
		return pipeShareData, fmt.Errorf("invalid pipe type share: %v", pipeShareData.Type)
	}

	if err != nil {
		return pipeShareData, err
	}
	return pipeShareData, nil
}

func (p pipeShareActions) CreatePipeReceiver(receiver models.SharedPipeReceiver) (models.SharedPipeReceiver, error) {
	query := `
			INSERT INTO shared_pipe_receivers (sharer_id, shared_pipe_id, receiver_id, is_accepted)
			VALUES ($1, $2, $3, $4)
			RETURNING id, sharer_id, shared_pipe_id, receiver_id, is_accepted, created_at, modified_at
	`
	err := p.Db.QueryRow(query, receiver.SharerId, receiver.SharedPipeId, receiver.ReceiverID, receiver.IsAccepted).Scan(
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

func (p pipeShareActions) GetSharedPipe(pipeId int64) (models.SharedPipe, error) {
	var sharedPipe models.SharedPipe
	query := "SELECT id, sharer_id, pipe_id, type, code, created_at FROM shared_pipes WHERE pipe_id=$1 LIMIT 1"
	err := p.Db.QueryRow(query, pipeId).Scan(
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
	return sharedPipe, nil
}

func (p pipeShareActions) GetSharedPipeByCode(code string) (models.SharedPipe, error) {
	var sharedPipe models.SharedPipe
	query := "SELECT id,sharer_id, pipe_id, type, code, created_at FROM shared_pipes WHERE code=$1 LIMIT 1"
	err := p.Db.QueryRow(query, code).Scan(
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

func (p pipeShareActions) GetReceivedPipeRecord(pipeId, userId int64) (models.SharedPipeReceiver, error) {
	logger := zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()

	var sharedPipe models.SharedPipeReceiver
	query := "SELECT id, shared_pipe_id, receiver_id FROM shared_pipe_receivers WHERE shared_pipe_id=$1 AND receiver_id=$2 LIMIT 1"
	err := p.Db.QueryRow(query, pipeId, userId).Scan(
		&sharedPipe.ID,
		&sharedPipe.SharedPipeId,
		&sharedPipe.ReceiverID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.SharedPipeReceiver{}, ErrNoRecord
		}
		logger.Info().Msg("There's a specific error somewhere")
		logger.Err(err)
		return models.SharedPipeReceiver{}, err
	}
	return sharedPipe, nil
}

func (p pipeShareActions) AcceptPrivateShare(receiver models.SharedPipeReceiver) (models.SharedPipeReceiver, error) {
	query := `
	UPDATE shared_pipe_receivers 
	SET 
	    is_accepted=true 
	WHERE receiver_id=$1 AND shared_pipe_id=$2
	RETURNING id, sharer_id, shared_pipe_id, receiver_id, is_accepted, created_at, modified_at
	`

	err := p.Db.QueryRow(query, receiver.ReceiverID, receiver.SharedPipeId).Scan(
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
