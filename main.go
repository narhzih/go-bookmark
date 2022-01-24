package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"gitlab.com/trencetech/mypipe-api/pkg/service/mailer"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	psh "github.com/platformsh/config-reader-go/v2"
	"github.com/rs/zerolog"
	"gitlab.com/trencetech/mypipe-api/db"
	"gitlab.com/trencetech/mypipe-api/pkg/api"
	"gitlab.com/trencetech/mypipe-api/pkg/service"
)

func main() {
	logger := zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()

	// ONly require .env file on local machine
	if onPlatformDotSh() == false {
		err := godotenv.Load(".env")
		if err != nil {
			logger.Err(err).Msg("Could not load environment variables")
			os.Exit(1)
		}
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

	mailerP := initMailer()
	apiService := service.NewService(db, jwtConfig, mailerP)
	apiHandler := api.NewHandler(apiService, logger)
	router := gin.Default()
	router.Use(cors.Default())
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
	var err error
	dsn, err := getSqlConnectionString(logger)
	if err != nil {
		return db.Database{}, err
	}
	return db.Connect(dsn, logger)
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

func initMailer() *mailer.Mailer {
	logger := zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()
	var mailerP *mailer.Mailer
	var config mailer.MailConfig
	config.Password = os.Getenv("MAIL_PASSWORD")
	config.Username = os.Getenv("MAIL_USERNAME")
	config.SmtpHost = os.Getenv("MAIL_HOST")
	config.SmtpPort = os.Getenv("MAIL_PORT")
	config.MailFrom = "My Pipe Desk <desk@mypipe.app>"
	logger.Info().Msg(config.Username + " is the username")
	logger.Info().Msg(config.Password + " is the password")
	mailerP = mailer.NewMailer(config)
	return mailerP
}

func getSqlConnectionString(logger zerolog.Logger) (string, error) {
	var postgresPort int
	var connectionString string
	var err error
	//platformShConfig, err := psh.NewRuntimeConfig()
	//if err != nil {
	//
	//}
	// This means we're on platform.sh
	//credentials, err := platformShConfig.Credentials("database")
	//if err != nil {
	//	logger.Err(err).Msg("Error occurred while trying to connect to platform.sh database")
	//	return "", err
	//}
	//
	//connectionString, err = sqldsn.FormattedCredentials(credentials)
	//if err != nil {
	//	logger.Err(err).Msg("Invalid formatted credentials")
	//	return "", nil
	//}

	// This means we're not on platform.sh
	// Format the connection string based
	// on variables in .env
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

func onPlatformDotSh() bool {
	_, err := psh.NewRuntimeConfig()
	if err != nil {
		return false
	}

	return true
}
