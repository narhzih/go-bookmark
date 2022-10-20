package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mypipeapp/mypipeapi/cmd/api/internal"
	"github.com/mypipeapp/mypipeapi/cmd/api/routes"
	"github.com/mypipeapp/mypipeapi/cmd/api/services"
	"github.com/mypipeapp/mypipeapi/db/actions/postgres"
	"github.com/mypipeapp/mypipeapi/db/repository"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func serveApp(db *sql.DB, logger zerolog.Logger) {

	repositories := repository.Repositories{
		User:                postgres.NewUserActions(db, logger),
		Pipe:                postgres.NewPipeActions(db, logger),
		PipeShare:           postgres.NewPipeShareActions(db, logger),
		Bookmark:            postgres.NewBookmarkActions(db, logger),
		Notification:        postgres.NewNotificationActions(db, logger),
		AccountVerification: postgres.NewAccountVerificationActions(db, logger),
	}

	jwtConfig, err := initJWTConfig()
	if err != nil {
		logger.Err(err).Msg("jwt config")
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	mailerP := initMailer()
	defer stop()

	// Create application instance
	app := internal.Application{
		Repositories: repositories,
		Logger:       logger,
		Services: services.Services{
			Repositories: repositories,
			Logger:       logger,
			JWTConfig:    jwtConfig,
			Mailer:       mailerP,
		},
	}

	// setup router
	router := gin.Default()
	router.Use(cors.Default())
	rg := router.Group("/v1")
	routes.BootRoutes(app, rg)

	// Start application server
	_, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		logger.Err(err).Msg("Unable to bind port")
	}

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

	// run the reset of the application and set up the routes
}
