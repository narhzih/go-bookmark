package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mypipeapp/mypipeapi/db/models"
	"github.com/mypipeapp/mypipeapi/db/repository"
	"github.com/rs/zerolog"
	"strings"
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
			Code:         pipeShareData.Code,
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

// CreatePipeReceiver creates a receiver for a pipe share
// receiver record for public pipe share is automatically marked as accepted
// while private pipe shares needs to be accepted by the receiver in another step
func (p pipeShareActions) CreatePipeReceiver(receiver models.SharedPipeReceiver) (models.SharedPipeReceiver, error) {
	query := `
	INSERT INTO shared_pipe_receivers 
	    (sharer_id, shared_pipe_id, receiver_id, code, is_accepted)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, sharer_id, shared_pipe_id, receiver_id, code, is_accepted, created_at, modified_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := p.Db.QueryRowContext(ctx, query,
		receiver.SharerId,
		receiver.SharedPipeId,
		receiver.ReceiverID,
		receiver.Code,
		receiver.IsAccepted,
	).Scan(
		&receiver.ID,
		&receiver.SharerId,
		&receiver.SharedPipeId,
		&receiver.ReceiverID,
		&receiver.Code,
		&receiver.IsAccepted,
		&receiver.CreatedAt,
		&receiver.ModifiedAt,
	)
	if err != nil {
		return receiver, err
	}
	return receiver, nil
}

// GetSharedPipe retrieves a pipe share record for a particular pipe
func (p pipeShareActions) GetSharedPipe(pipeId int64, shareType string) (models.SharedPipe, error) {
	if strings.TrimSpace(shareType) == "" {
		shareType = models.PipeShareTypePrivate
	}
	var sharedPipe models.SharedPipe
	query := `
	SELECT id, sharer_id, pipe_id, type, code, created_at 
	FROM shared_pipes 
	WHERE pipe_id=$1 AND type=$2
	LIMIT 1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := p.Db.QueryRowContext(ctx, query, pipeId, shareType).Scan(
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

// GetSharedPipeByCode retrieves a shared pipe record by share code
func (p pipeShareActions) GetSharedPipeByCode(code string) (models.SharedPipe, error) {
	var sharedPipe models.SharedPipe
	query := `
	SELECT id,sharer_id, pipe_id, type, code, created_at 
	FROM shared_pipes 
	WHERE code=$1 
	LIMIT 1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := p.Db.QueryRowContext(ctx, query, code).Scan(
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

// GetReceivedPipeRecord retrieves a received pipe record by pipe id and a designated user id
func (p pipeShareActions) GetReceivedPipeRecord(pipeId, userId int64) (models.SharedPipeReceiver, error) {

	var sharedPipe models.SharedPipeReceiver
	query := `
	SELECT id, sharer_id, shared_pipe_id, receiver_id, code, is_accepted
	FROM shared_pipe_receivers 
	WHERE shared_pipe_id=$1 AND receiver_id=$2 
	LIMIT 1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := p.Db.QueryRowContext(ctx, query, pipeId, userId).Scan(
		&sharedPipe.ID,
		&sharedPipe.SharerId,
		&sharedPipe.SharedPipeId,
		&sharedPipe.ReceiverID,
		&sharedPipe.Code,
		&sharedPipe.IsAccepted,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.SharedPipeReceiver{}, ErrNoRecord
		}
		p.Logger.Err(err)
		return models.SharedPipeReceiver{}, err
	}
	return sharedPipe, nil
}

// GetReceivedPipeRecordByCode retrieves a received pipe record by pipe id and a designated user id
func (p pipeShareActions) GetReceivedPipeRecordByCode(code string, userId int64) (models.SharedPipeReceiver, error) {

	var sharedPipe models.SharedPipeReceiver
	query := `
	SELECT id, sharer_id, shared_pipe_id, receiver_id, code, is_accepted
	FROM shared_pipe_receivers 
	WHERE code=$1 AND receiver_id=$2 
	LIMIT 1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := p.Db.QueryRowContext(ctx, query, code, userId).Scan(
		&sharedPipe.ID,
		&sharedPipe.SharerId,
		&sharedPipe.SharedPipeId,
		&sharedPipe.ReceiverID,
		&sharedPipe.Code,
		&sharedPipe.IsAccepted,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.SharedPipeReceiver{}, ErrNoRecord
		}
		p.Logger.Err(err)
		return models.SharedPipeReceiver{}, err
	}
	return sharedPipe, nil
}

// AcceptPrivateShare accept a private share
func (p pipeShareActions) AcceptPrivateShare(receiver models.SharedPipeReceiver) (models.SharedPipeReceiver, error) {
	query := `
	UPDATE shared_pipe_receivers 
	SET 
	    is_accepted=true, modified_at=now()
	WHERE id=$1
	RETURNING id, sharer_id, shared_pipe_id, receiver_id, code, is_accepted, created_at, modified_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := p.Db.QueryRowContext(ctx, query, receiver.ID).Scan(
		&receiver.ID,
		&receiver.SharerId,
		&receiver.SharedPipeId,
		&receiver.ReceiverID,
		&receiver.Code,
		&receiver.IsAccepted,
		&receiver.CreatedAt,
		&receiver.ModifiedAt,
	)

	if err != nil {
		return receiver, err
	}
	return receiver, nil
}
