package main

import (
	"database/sql"
	"fmt"
	"github.com/mypipeapp/mypipeapi/db"
	"github.com/rs/zerolog"
	"os"
	"strconv"
)

func initDb(logger zerolog.Logger) (*sql.DB, error) {
	var err error
	dsn, err := getSqlConnectionString(logger)
	if err != nil {
		return nil, err
	}
	dbInit, _ := db.Connect(dsn, logger)
	return dbInit.Conn, nil
}

func getSqlConnectionString(logger zerolog.Logger) (string, error) {
	var postgresPort int
	var connectionString string
	var err error
	postgresPort, err = strconv.Atoi(os.Getenv("POSTGRES_DB_PORT"))
	if err != nil {
		logger.Err(err).Msg("Error coming from parsing DB_PORT")
		return "", err
	}
	dbConfig := db.Config{
		Host:           os.Getenv("POSTGRES_DB_HOST"),
		Port:           postgresPort,
		DbName:         os.Getenv("POSTGRES_DB"),
		Username:       os.Getenv("POSTGRES_USER"),
		Password:       os.Getenv("POSTGRES_PASSWORD"),
		ConnectionMode: os.Getenv("DB_SSL_MODE"),
		Logger:         logger,
	}

	connectionString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.DbName, dbConfig.ConnectionMode)

	return connectionString, nil
}
