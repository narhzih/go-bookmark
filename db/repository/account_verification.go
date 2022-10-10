package repository

import "gitlab.com/trencetech/mypipe-api/db/models"

type AccountVerificationRepository interface {
	CreateVerification(accountVerification models.AccountVerification) (models.AccountVerification, error)
	GetAccountVerificationByToken(token string) (models.AccountVerification, error)
	DeleteVerification(token string) (bool, error)
}
