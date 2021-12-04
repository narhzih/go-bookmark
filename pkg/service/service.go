package service

import (
	"fmt"
	"gitlab.com/trencetech/mypipe-api/db"
)

// Service More fields will be added to the service struct later in the future
var (
	ErrFileTooLarge = fmt.Errorf("file too large")
)

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
