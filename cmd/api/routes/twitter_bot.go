package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mypipeapp/mypipeapi/cmd/api/handlers"
	"github.com/mypipeapp/mypipeapi/cmd/api/internal"
)

func setupTwitterBotRoutes(app internal.Application, routeGroup *gin.RouterGroup) {
	h := handlers.NewTwitterBotHandler(app)
	bot := routeGroup.Group("/bot")
	bot.POST("/twitter/add-to-pipe", h.BotAddToPipe)
}
