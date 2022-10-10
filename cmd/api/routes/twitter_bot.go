package routes

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/trencetech/mypipe-api/cmd/api/handlers"
	"gitlab.com/trencetech/mypipe-api/cmd/api/internal"
)

func setupTwitterBotRoutes(app internal.Application, routeGroup *gin.RouterGroup) {
	h := handlers.NewTwitterBotHandler(app)
	bot := routeGroup.Group("/bot")
	bot.POST("/twitter/add-to-pipe", h.BotAddToPipe)
}
