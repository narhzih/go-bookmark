package repository

import "github.com/mypipeapp/mypipeapi/db/models"

type AccountVerificationRepository interface {
	CreateVerification(accountVerification models.AccountVerification) (models.AccountVerification, error)
	GetAccountVerificationByToken(token string) (models.AccountVerification, error)
	DeleteVerification(token string) (bool, error)
}
