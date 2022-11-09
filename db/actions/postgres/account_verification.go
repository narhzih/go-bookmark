package postgres

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"github.com/mypipeapp/mypipeapi/db/models"
	"github.com/mypipeapp/mypipeapi/db/repository"
	"github.com/rs/zerolog"
	"time"
)

type accountVerificationActions struct {
	Db     *sql.DB
	Logger zerolog.Logger
}

func NewAccountVerificationActions(db *sql.DB, logger zerolog.Logger) repository.AccountVerificationRepository {
	return accountVerificationActions{
		Db:     db,
		Logger: logger,
	}
}

// CreateVerification creates a verification token record for a user
func (a accountVerificationActions) CreateVerification(accountVerification models.AccountVerification) (models.AccountVerification, error) {
	query := `
	INSERT INTO account_verifications (user_id, token, expires_at) 
	VALUES ($1, $2, $3) 
	RETURNING id, user_id, token, created_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := a.Db.QueryRowContext(
		ctx,
		query,
		accountVerification.UserID,
		accountVerification.Token,
		accountVerification.ExpiresAt,
	).Scan(
		&accountVerification.ID,
		&accountVerification.UserID,
		&accountVerification.Token,
		&accountVerification.CreatedAt,
	)
	if err != nil {
		if dbErr, ok := err.(*pq.Error); ok {
			if dbErr.Code == "23505" {
				return models.AccountVerification{}, ErrRecordExists
			}
		}
		return models.AccountVerification{}, err
	}
	return accountVerification, nil
}

// GetAccountVerificationByToken fetches verification record by token
func (a accountVerificationActions) GetAccountVerificationByToken(token string) (models.AccountVerification, error) {
	var accountVerification models.AccountVerification
	query := `
	SELECT id, user_id, used, token, created_at 
	FROM account_verifications 
	WHERE token=$1 LIMIT 1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := a.Db.QueryRowContext(ctx, query, token).Scan(
		&accountVerification.ID,
		&accountVerification.UserID,
		&accountVerification.Used,
		&accountVerification.Token,
		&accountVerification.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return models.AccountVerification{}, ErrNoRecord
		}
		return models.AccountVerification{}, nil
	}

	return accountVerification, nil
}

// DeleteVerification completely removes a verification record of a user
func (a accountVerificationActions) DeleteVerification(token string) (bool, error) {
	query := `DELETE FROM account_verifications WHERE token=$1`
	_, err := a.Db.Exec(query, token)
	if err != nil {
		return false, err
	}
	return true, nil
}
