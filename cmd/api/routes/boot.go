package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mypipeapp/mypipeapi/cmd/api/internal"
	"net/http"
)

func BootRoutes(app internal.Application, routeGroup *gin.RouterGroup) {
	// register healthz route
	routeGroup.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to mypipeapp api",
		})
	})

	// setup necessary routes
	setupAuthRoutes(app, routeGroup)
	setupUserRoutes(app, routeGroup)
	setupPipeRoutes(app, routeGroup)
	setupNotificationRoutes(app, routeGroup)
	setupTwitterBotRoutes(app, routeGroup)
	setupParserRoutes(app, routeGroup)
	setupSearchRoutes(app, routeGroup)
}
