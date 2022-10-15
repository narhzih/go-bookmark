package routes

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/trencetech/mypipe-api/cmd/api/handlers"
	"gitlab.com/trencetech/mypipe-api/cmd/api/internal"
	"gitlab.com/trencetech/mypipe-api/cmd/api/middlewares"
)

func setupAuthRoutes(app internal.Application, routeGroup *gin.RouterGroup) {
	h := handlers.NewAuthHandler(app)
	routeGroup.POST("/sign-up", h.EmailSignUp)
	routeGroup.POST("/sign-in", h.EmailLogin)
	routeGroup.POST("/verify-account/:token", h.VerifyAccount)
	routeGroup.POST("/forgot-password", h.ForgotPassword)
	routeGroup.POST("/verify-reset-token/:token", h.VerifyPasswordResetToken)
	routeGroup.POST("/reset-password/:token", h.ResetPassword)

	routeGroup.POST("/google-auth", h.SignInWithGoogle)

	authApi := routeGroup.Group("/auth")
	authApi.Use(middlewares.AuthRequired(app, app.Services.JWTConfig.Key, app.Logger))
	authApi.POST("/twitter/connect-account", h.ConnectTwitterAccount)
}
