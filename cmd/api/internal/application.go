package internal

import (
	"gitlab.com/trencetech/mypipe-api/db/repository"
	"gitlab.com/trencetech/mypipe-api/pkg/service"
)

type Application struct {
	Repositories repository.Repositories
	Services     service.Service
	Helpers      string
}
