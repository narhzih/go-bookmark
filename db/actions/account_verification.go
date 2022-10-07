package actions

import (
	"database/sql"
	"github.com/rs/zerolog"
	"gitlab.com/trencetech/mypipe-api/db/models"
	"gitlab.com/trencetech/mypipe-api/db/repository"
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

func (a accountVerificationActions) CreateVerification(accountVerification models.AccountVerification) (models.AccountVerification, error) {
	query := `
				INSERT INTO account_verifications (user_id, token, expires_at) 
				VALUES ($1, $2, $3) 
				RETURNING id, user_id, token, created_at
			`
	err := a.Db.QueryRow(query, accountVerification.UserID, accountVerification.Token, accountVerification.ExpiresAt).Scan(
		&accountVerification.ID,
		&accountVerification.UserID,
		&accountVerification.Token,
		&accountVerification.CreatedAt,
	)
	if err != nil {
		return models.AccountVerification{}, err
	}
	return accountVerification, nil
}

func (a accountVerificationActions) GetAccountVerificationByToken(token string) (models.AccountVerification, error) {
	var accountVerification models.AccountVerification
	query := `SELECT id, user_id, used, token, created_at FROM account_verifications WHERE token=$1 LIMIT 1`
	if err := a.Db.QueryRow(query, token).Scan(
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

func (a accountVerificationActions) DeleteVerification(token string) (bool, error) {
	query := "DELETE FROM account_verifications WHERE token=$1"
	_, err := a.Db.Exec(query, token)
	if err != nil {
		return false, err
	}
	return true, nil
}
