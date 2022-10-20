package internal

import (
	"github.com/mypipeapp/mypipeapi/cmd/api/services"
	"github.com/mypipeapp/mypipeapi/db/repository"
	"github.com/rs/zerolog"
)

type Application struct {
	Repositories repository.Repositories
	Services     services.Services
	Logger       zerolog.Logger
}
