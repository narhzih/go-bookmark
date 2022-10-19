package main

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	psh "github.com/platformsh/config-reader-go/v2"
	"github.com/rs/zerolog"
	"gitlab.com/trencetech/mypipe-api/cmd/api/services"
	"gitlab.com/trencetech/mypipe-api/cmd/api/services/mailer"
	"os"
	"strconv"
)

func main() {
	logger := zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()

	// set application environment variable loader
	appEnv := os.Getenv("APP_ENV")
	var err error
	if appEnv == "dev" {
		logger.Info().Msg("Loading prod env")
		err = godotenv.Load(".env")
	}
	//else {
	//	logger.Info().Msg("Loading dev env")
	//	err = godotenv.Load(".env")
	//}
	if err != nil {
		logger.Err(err).Msg("Could not load environment variables")
		os.Exit(1)
	}

	db, err := initDb(logger)
	if err != nil {
		logger.Err(err).Msg("An error occurred")
		os.Exit(1)
	}
	logger.Info().Msg("Established connection with api database")

	serveApp(db, logger)

	//
	//go func() {
	//	if err := srv.ListenAndServe(); err != nil {
	//		logger.Err(err).Msg("listen")
	//	}
	//}()
	//
	//<-ctx.Done()
	//stop()
	//
	//logger.Info().Msg("shutting down gracefully, press Ctrl+C again to force")
	//
	//// The context is used to inform the server it has 5 seconds to finish
	//// the request it is currently handling
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//if err := srv.Shutdown(ctx); err != nil {
	//	logger.Fatal().Msg(fmt.Sprintf("Server forced to shutdown: %s", err))
	//}
	//logger.Info().Msg("exiting server")
}

func initJWTConfig() (services.JWTConfig, error) {
	var expiresIn int
	var key string
	var err error
	var cfg services.JWTConfig

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

func onPlatformDotSh() bool {
	_, err := psh.NewRuntimeConfig()
	if err != nil {
		return false
	}

	return true
}
