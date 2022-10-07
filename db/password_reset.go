package db

import (
	"database/sql"
	"gitlab.com/trencetech/mypipe-api/db/models"
)

func (db Database) CreatePasswordResetRecord(user models.User, token string) (models.PasswordReset, error) {
	var passwordReset models.PasswordReset
	query := `INSERT INTO password_resets (user_id, token) VALUES ($1, $2) RETURNING id, user_id, token, created_at`
	err := db.Conn.QueryRow(query, user.ID, token).Scan(
		&passwordReset.ID,
		&passwordReset.UserID,
		&passwordReset.Token,
		&passwordReset.CreatedAt,
	)
	if err != nil {
		return passwordReset, err
	}
	return passwordReset, nil
}

func (db Database) GetPasswordResetRecord(token string) (models.PasswordReset, error) {
	var passwordReset models.PasswordReset
	query := `SELECT id, user_id, token, created_at, validated FROM password_resets WHERE token=$1 LIMIT 1`
	err := db.Conn.QueryRow(query, token).Scan(
		&passwordReset.ID,
		&passwordReset.UserID,
		&passwordReset.Token,
		&passwordReset.CreatedAt,
		&passwordReset.Validated,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return passwordReset, ErrNoRecord
		}
		return passwordReset, err
	}
	return passwordReset, nil
}

func (db Database) UpdatePasswordResetRecord(token string) (models.PasswordReset, error) {
	var passwordReset models.PasswordReset
	query := `UPDATE password_resets SET validated=true WHERE token=$1 RETURNING id, user_id, token, created_at, validated`
	err := db.Conn.QueryRow(query, token).Scan(
		&passwordReset.ID,
		&passwordReset.UserID,
		&passwordReset.Token,
		&passwordReset.CreatedAt,
		&passwordReset.Validated,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return passwordReset, ErrNoRecord
		}
		return passwordReset, err
	}
	return passwordReset, nil
}

func (db Database) DeletePasswordResetRecord(token string) error {
	query := `DELETE FROM password_resets WHERE token=$1`
	_, err := db.Conn.Exec(query, token)
	if err != nil {
		return err
	}
	return nil
}
