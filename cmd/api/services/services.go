package services

import (
	"github.com/mypipeapp/mypipeapi/cmd/api/services/mailer"
	"github.com/mypipeapp/mypipeapi/db/repository"
	"github.com/rs/zerolog"
)

type Services struct {
	Mailer       *mailer.Mailer
	Repositories repository.Repositories
	Logger       zerolog.Logger
	JWTConfig    JWTConfig
}
