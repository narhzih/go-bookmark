package internal

import (
	"github.com/rs/zerolog"
	"gitlab.com/trencetech/mypipe-api/cmd/api/services"
	"gitlab.com/trencetech/mypipe-api/db/repository"
)

type Application struct {
	Repositories repository.Repositories
	Services     services.Services
	Logger       zerolog.Logger
}
