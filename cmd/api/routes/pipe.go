package routes

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/trencetech/mypipe-api/cmd/api/handlers"
	"gitlab.com/trencetech/mypipe-api/cmd/api/internal"
)

func setupPipeRoutes(app internal.Application, routeGroup *gin.RouterGroup) {
	h := handlers.NewPipeHandler(app)
	bookmarkH := handlers.NewBookmarkHandler(app)
	pipeShareH := handlers.NewPipeShareHandler(app)

	pipe := routeGroup.Group("/pipe")
	pipe.Use(handlers.AuthRequired(app.Services.JWTConfig.Key, app.Logger))
	pipe.POST("/", h.CreatePipe)
	pipe.GET("/:id", h.GetPipe)
	pipe.POST("/:id/share", pipeShareH.SharePipe)
	pipe.PUT("/:id", h.UpdatePipe)
	pipe.DELETE("/:id", h.DeletePipe)
	pipe.GET("/all", h.GetPipes)
	pipe.GET("/all/steroids", h.GetPipeWithResource)
	pipe.GET("/preview", pipeShareH.PreviewPipe)
	pipe.POST("/add-pipe", pipeShareH.AddPipe)

	pipe.GET("/:id/bookmarks", bookmarkH.GetBookmarks)
	pipe.POST("/:id/bookmark", bookmarkH.CreateBookmark)
	pipe.GET("/:id/bookmark/:bmId", bookmarkH.GetBookmark)
	pipe.DELETE("/:id/bookmark/:bmId", bookmarkH.DeleteBookmark)
}
