package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mypipeapp/mypipeapi/cmd/api/handlers"
	"github.com/mypipeapp/mypipeapi/cmd/api/internal"
	"github.com/mypipeapp/mypipeapi/cmd/api/middlewares"
)

func setupUserRoutes(app internal.Application, routeGroup *gin.RouterGroup) {
	h := handlers.NewUserHandler(app)

	user := routeGroup.Group("/user")
	user.Use(middlewares.AuthRequired(app, app.Services.JWTConfig.Key))
	user.GET("/profile", h.UserProfile)
	user.PATCH("/profile", h.EditProfile)
	user.PATCH("/profile/change-password", h.ChangePassword)
	user.POST("/profile/cover-photo", h.UploadCoverPhoto)
}
