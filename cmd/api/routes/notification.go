package routes

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/trencetech/mypipe-api/cmd/api/handlers"
	"gitlab.com/trencetech/mypipe-api/cmd/api/internal"
)

func setupNotificationRoutes(app internal.Application, routeGroup *gin.RouterGroup) {
	h := handlers.NewNotificationHandler(app)
	notification := routeGroup.Group("/notifications")
	notification.Use(handlers.AuthRequired(app.Services.JWTConfig.Key, app.Logger))
	notification.GET("/", h.GetNotifications)
	notification.POST("/update-device-tokens", h.UpdateUserDeviceTokens)
	notification.GET("/:notificationId", h.GetNotification)
}
