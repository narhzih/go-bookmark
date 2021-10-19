package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"gitlab.com/gowagr/mypipe-api/db"
)

func main() {
	godotenv.Load(".env")
	logger := zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()
	_, err := initDb(logger)
	if err != nil {
		logger.Err(err).Msg("An error occurred")
		os.Exit(1)
	}
	logger.Info().Msg("Established connection with api database")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := gin.Default()
	rg := router.Group("/v1")
	rg.GET("/test", FirstCall)
	addr := fmt.Sprintf(":%d", 5555)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Err(err).Msg("listen")
		}
	}()

	<-ctx.Done()
}

func initDb(logger zerolog.Logger) (db.Database, error) {
	var postgresPort int
	var err error
	postgresPort, err = strconv.Atoi(os.Getenv("POSTGRES_DB_PORT"))
	if err != nil {
		return db.Database{}, err
	}
	dbConfig := db.Config{
		Host:     os.Getenv("POSTGRES_DB_HOST"),
		Port:     postgresPort,
		DbName:   os.Getenv("POSTGRES_DB"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Logger:   logger,
	}

	return db.Connect(dbConfig)
}

func FirstCall(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Connection successful",
	})

}
