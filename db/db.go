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
	Host     string
	Username string
	Password string
	Port     int
	DbName   string
	Logger   zerolog.Logger
}

//<--- User and user auth structs

func Connect(config Config) (Database, error) {
	db := Database{}
	dsn := fmt.Sprintf("host=%s port=%d username=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, config.DbName)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return db, err
	}
	db.Conn = conn
	db.Logger = config.Logger
	err = db.Conn.Ping()
	if err != nil {
		return db, err
	}

	return db, nil
}
