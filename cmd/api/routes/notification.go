package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mypipeapp/mypipeapi/cmd/api/handlers"
	"github.com/mypipeapp/mypipeapi/cmd/api/internal"
	"github.com/mypipeapp/mypipeapi/cmd/api/middlewares"
)

func setupNotificationRoutes(app internal.Application, routeGroup *gin.RouterGroup) {
	h := handlers.NewNotificationHandler(app)
	notification := routeGroup.Group("/notifications")
	notification.Use(middlewares.AuthRequired(app, app.Services.JWTConfig.Key))
	notification.GET("/", h.GetNotifications)
	notification.POST("/update-device-tokens", h.UpdateUserDeviceTokens)
	notification.GET("/:notificationId", h.GetNotification)
}
