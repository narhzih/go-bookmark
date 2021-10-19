package service

import (
	"gitlab.com/gowagr/mypipe-api/db"
)

// More fields will be added to the service struct later in the future
type Service struct {
	DB db.Database
}

func NewService(dbHandle db.Database) Service {
	return Service{
		DB: dbHandle,
	}
}
