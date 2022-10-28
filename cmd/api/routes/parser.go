package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mypipeapp/mypipeapi/cmd/api/handlers"
	"github.com/mypipeapp/mypipeapi/cmd/api/internal"
)

func setupParserRoutes(app internal.Application, routeGroup *gin.RouterGroup) {
	h := handlers.NewParserHandler(app)
	parser := routeGroup.Group("/parse-link")
	//parser.Use(AuthRequired(h.service.JWTConfig.Key, h.logger))
	parser.POST("/twitter", h.TwitterLinkParser)
	parser.POST("/twitter/thread", h.GetCompleteThreadOfATweet)
	parser.POST("/youtube", h.YoutubeLinkParser)
	parser.POST("/others", h.ParseLink)
}
