package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mypipeapp/mypipeapi/cmd/api/handlers"
	"github.com/mypipeapp/mypipeapi/cmd/api/internal"
	"github.com/mypipeapp/mypipeapi/cmd/api/middlewares"
)

func setupPipeRoutes(app internal.Application, routeGroup *gin.RouterGroup) {
	h := handlers.NewPipeHandler(app)
	bookmarkH := handlers.NewBookmarkHandler(app)
	pipeShareH := handlers.NewPipeShareHandler(app)

	pipe := routeGroup.Group("/pipe")
	pipe.Use(middlewares.AuthRequired(app, app.Services.JWTConfig.Key))

	pipe.POST("/", h.CreatePipe)
	pipe.GET("/:id", h.GetPipe)
	pipe.POST("/bookmark", bookmarkH.CreateBookmark)
	pipe.POST("/:id/share", pipeShareH.SharePipe)
	pipe.PUT("/:id", h.UpdatePipe)
	pipe.DELETE("/:id", h.DeletePipe)
	pipe.GET("/all", h.GetPipes)
	pipe.GET("/preview", pipeShareH.PreviewPipe)
	pipe.POST("/add-pipe", pipeShareH.AddPipe)

	pipe.GET("/:id/bookmarks", bookmarkH.GetBookmarks)
	pipe.GET("/:id/bookmark/:bmId", bookmarkH.GetBookmark)
	pipe.DELETE("/:id/bookmark/:bmId", bookmarkH.DeleteBookmark)
}
