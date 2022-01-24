package db

import (
	"database/sql"
	"fmt"

	"github.com/rs/zerolog"
)

var (
	ErrRecordExists = fmt.Errorf("row with the same value already exits")
	ErrNoRecord     = fmt.Errorf("no matching row was found")
)

type Database struct {
	Conn   *sql.DB
	Logger zerolog.Logger
}

type Config struct {
	Host           string
	Username       string
	Password       string
	Port           int
	DbName         string
	ConnectionMode string
	Logger         zerolog.Logger
}

//<--- User and user auth structs

func Connect(connectionString string, logger zerolog.Logger) (Database, error) {
	db := Database{}
	conn, err := sql.Open("postgres", connectionString)
	if err != nil {
		logger.Err(err).Msg("Error was coming from database initialization")
		return db, err
	}
	db.Conn = conn
	db.Logger = logger
	err = db.Conn.Ping()
	if err != nil {
		logger.Err(err).Msg("Cannot ping database because error occurred while pinging")
		return db, err
	}

	return db, nil
}
