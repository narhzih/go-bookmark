package repository

import "github.com/mypipeapp/mypipeapi/db/models"

type PasswordResetRepository interface {
	CreatePasswordResetRecord(user models.User, token string) (models.PasswordReset, error)
	GetPasswordResetRecord(token string) (models.PasswordReset, error)
	UpdatePasswordResetRecord(token string) (models.PasswordReset, error)
	DeletePasswordResetRecord(token string) error
}
