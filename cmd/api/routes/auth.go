package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mypipeapp/mypipeapi/cmd/api/handlers"
	"github.com/mypipeapp/mypipeapi/cmd/api/internal"
	"github.com/mypipeapp/mypipeapi/cmd/api/middlewares"
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
	authApi.Use(middlewares.AuthRequired(app, app.Services.JWTConfig.Key))
	authApi.POST("/twitter/connect-account", h.ConnectTwitterAccount)
}
