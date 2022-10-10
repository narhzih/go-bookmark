package routes

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/trencetech/mypipe-api/cmd/api/handlers"
	"gitlab.com/trencetech/mypipe-api/cmd/api/internal"
)

func setupParserRoutes(app internal.Application, routeGroup *gin.RouterGroup) {
	h := handlers.NewParserHandler(app)
	parser := routeGroup.Group("/parse-link")
	//parser.Use(AuthRequired(h.service.JWTConfig.Key, h.logger))
	parser.POST("/twitter", h.TwitterLinkParser)
	parser.POST("/youtube", h.YoutubeLinkParser)
	parser.POST("/others", h.ParseLink)
}
