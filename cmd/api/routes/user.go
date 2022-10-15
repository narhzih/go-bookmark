package routes

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/trencetech/mypipe-api/cmd/api/handlers"
	"gitlab.com/trencetech/mypipe-api/cmd/api/internal"
	"gitlab.com/trencetech/mypipe-api/cmd/api/middlewares"
)

func setupUserRoutes(app internal.Application, routeGroup *gin.RouterGroup) {
	h := handlers.NewUserHandler(app)

	user := routeGroup.Group("/user")
	user.Use(middlewares.AuthRequired(app, app.Services.JWTConfig.Key, app.Logger))
	user.GET("/profile", h.UserProfile)
	user.PATCH("/profile", h.EditProfile)
	user.PATCH("/profile/change-password", h.ChangePassword)
	user.POST("/profile/cover-photo", h.UploadCoverPhoto)
}
