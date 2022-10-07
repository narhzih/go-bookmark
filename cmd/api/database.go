package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"gitlab.com/trencetech/mypipe-api/db"
	"os"
	"strconv"
)

func initDb(logger zerolog.Logger) (db.Database, error) {
	var err error
	dsn, err := getSqlConnectionString(logger)
	if err != nil {
		return db.Database{}, err
	}
	return db.Connect(dsn, logger)
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
