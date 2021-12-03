package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"gitlab.com/trencetech/mypipe-api/db"
	"gitlab.com/trencetech/mypipe-api/pkg/api"
	"gitlab.com/trencetech/mypipe-api/pkg/service"
)

func main() {
	logger := zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()

	// ONly require .env file on local machine
	if os.Getenv("PORT") == "" || len(os.Getenv("PORT")) <= 0 {
		logger.Info().Msg("Loading .env file")
		godotenv.Load(".env")
	}
	db, err := initDb(logger)
	if err != nil {
		logger.Err(err).Msg("An error occurred")
		os.Exit(1)
	}
	logger.Info().Msg("Established connection with api database")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	jwtConfig, err := initJWTConfig()
	if err != nil {
		logger.Err(err).Msg("jwt config")
	}

	apiService := service.NewService(db, jwtConfig)
	apiHandler := api.NewHandler(apiService, logger)
	router := gin.Default()
	rg := router.Group("/v1")
	apiHandler.Register(rg)
	appPort, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		logger.Err(err).Msg("Unable to bind port")
	}
	addr := fmt.Sprintf(":%d", appPort)
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
	stop()

	logger.Info().Msg("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal().Msg(fmt.Sprintf("Server forced to shutdown: %s", err))
	}
	logger.Info().Msg("exiting server")
}

func initDb(logger zerolog.Logger) (db.Database, error) {
	var postgresPort int
	var err error
	postgresPort, err = strconv.Atoi(os.Getenv("POSTGRES_DB_PORT"))
	if err != nil {
		logger.Err(err).Msg("Error coming from parsing DB_PORT")
		return db.Database{}, err
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

	return db.Connect(dbConfig)
}

func initJWTConfig() (service.JWTConfig, error) {
	var expiresIn int
	var key string
	var err error
	var cfg service.JWTConfig

	expiresIn, err = strconv.Atoi(os.Getenv("JWT_EXPIRES_IN"))
	if err != nil {
		return cfg, err
	}

	key = os.Getenv("JWT_SECRET")

	// enforce minimum length for JWT secret
	if len(key) < 64 {
		return cfg, fmt.Errorf("JWT_SECRET too short")
	}

	cfg.ExpiresIn = expiresIn
	cfg.Key = key
	cfg.Algo = jwt.SigningMethodHS256

	return cfg, nil
}
