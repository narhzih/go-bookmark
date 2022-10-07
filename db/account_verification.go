package db

import (
	"database/sql"
	"gitlab.com/trencetech/mypipe-api/db/models"
)

func (db Database) CreateVerification(accountVerification models.AccountVerification) (models.AccountVerification, error) {
	query := `
				INSERT INTO account_verifications (user_id, token, expires_at) 
				VALUES ($1, $2, $3) 
				RETURNING id, user_id, token, created_at
			`
	err := db.Conn.QueryRow(query, accountVerification.UserID, accountVerification.Token, accountVerification.ExpiresAt).Scan(
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

func (db Database) GetAccountVerificationByToken(token string) (models.AccountVerification, error) {
	var accountVerification models.AccountVerification
	query := `SELECT id, user_id, used, token, created_at FROM account_verifications WHERE token=$1 LIMIT 1`
	if err := db.Conn.QueryRow(query, token).Scan(
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

func (db Database) DeleteVerification(token string) (bool, error) {
	query := "DELETE FROM account_verifications WHERE token=$1"
	_, err := db.Conn.Exec(query, token)
	if err != nil {
		return false, err
	}
	return true, nil
}
