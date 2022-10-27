package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mypipeapp/mypipeapi/cmd/api/handlers"
	"github.com/mypipeapp/mypipeapi/cmd/api/internal"
	"github.com/mypipeapp/mypipeapi/cmd/api/middlewares"
)

func setupSearchRoutes(app internal.Application, routeGroup *gin.RouterGroup) {
	handler := handlers.NewSearchHandler(app)
	search := routeGroup.Group("search")
	search.Use(middlewares.AuthRequired(app, app.Services.JWTConfig.Key))
	search.GET("/", handler.Search)
}
