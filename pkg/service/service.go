package service

import (
	"gitlab.com/gowagr/mypipe-api/db"
)

// More fields will be added to the service struct later in the future
type Service struct {
	DB        db.Database
	JWTConfig JWTConfig
}

func NewService(dbHandle db.Database, jwtConfig JWTConfig) Service {
	return Service{
		DB:        dbHandle,
		JWTConfig: jwtConfig,
	}
}
