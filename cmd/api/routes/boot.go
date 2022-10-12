package routes

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/trencetech/mypipe-api/cmd/api/internal"
)

func BootRoutes(app internal.Application, routeGroup *gin.RouterGroup) {
	// setup necessary routes
	setupAuthRoutes(app, routeGroup)
	setupUserRoutes(app, routeGroup)
	setupPipeRoutes(app, routeGroup)
	setupNotificationRoutes(app, routeGroup)
	setupTwitterBotRoutes(app, routeGroup)
	setupParserRoutes(app, routeGroup)
}
