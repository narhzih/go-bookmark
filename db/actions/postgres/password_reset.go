package postgres

import (
	"database/sql"
	"github.com/rs/zerolog"
	"gitlab.com/trencetech/mypipe-api/db/models"
	"gitlab.com/trencetech/mypipe-api/db/repository"
)

type passwordResetActions struct {
	Db     *sql.DB
	Logger zerolog.Logger
}

func NewPasswordResetActions(db *sql.DB, logger zerolog.Logger) repository.PasswordResetRepository {
	return passwordResetActions{
		Db:     db,
		Logger: logger,
	}
}

func (p passwordResetActions) CreatePasswordResetRecord(user models.User, token string) (models.PasswordReset, error) {
	var passwordReset models.PasswordReset
	query := `INSERT INTO password_resets (user_id, token) VALUES ($1, $2) RETURNING id, user_id, token, created_at`
	err := p.Db.QueryRow(query, user.ID, token).Scan(
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

func (p passwordResetActions) GetPasswordResetRecord(token string) (models.PasswordReset, error) {
	var passwordReset models.PasswordReset
	query := `SELECT id, user_id, token, created_at, validated FROM password_resets WHERE token=$1 LIMIT 1`
	err := p.Db.QueryRow(query, token).Scan(
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

func (p passwordResetActions) UpdatePasswordResetRecord(token string) (models.PasswordReset, error) {
	var passwordReset models.PasswordReset
	query := `UPDATE password_resets SET validated=true WHERE token=$1 RETURNING id, user_id, token, created_at, validated`
	err := p.Db.QueryRow(query, token).Scan(
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

func (p passwordResetActions) DeletePasswordResetRecord(token string) error {
	query := `DELETE FROM password_resets WHERE token=$1`
	_, err := p.Db.Exec(query, token)
	if err != nil {
		return err
	}
	return nil
}
