package services

import (
	"github.com/rs/zerolog"
	"gitlab.com/trencetech/mypipe-api/cmd/api/services/mailer"
	"gitlab.com/trencetech/mypipe-api/db/repository"
)

type Services struct {
	Mailer       *mailer.Mailer
	Repositories repository.Repositories
	Logger       zerolog.Logger
	JWTConfig    JWTConfig
}
