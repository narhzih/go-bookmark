package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gitlab.com/gowagr/mypipe-api/pkg/service"
)

type Handler struct {
	service service.Service
	logger  zerolog.Logger
}

func NewHandler(service service.Service, logger zerolog.Logger) Handler {
	return Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) Register(routeGroup *gin.RouterGroup) {
	// This is where all routes will be registered
	routeGroup.GET("/test-route", TestCaller)
	routeGroup.POST("/google/signup", h.SingUpWithGoogle)
	routeGroup.POST("/google/singin", h.SignInWithGoogle)

	authApi := routeGroup.Group("/auth")
	authApi.Use(AuthRequired(h.service.JWTConfig.Key, h.logger))
	authApi.GET("/auth-test", TestCaller)

	pipe := routeGroup.Group("/pipe")
	pipe.Use(AuthRequired(h.service.JWTConfig.Key, h.logger))
	pipe.POST("/", h.CreatePipe)
	pipe.GET("/:id", h.GetPipe)
	pipe.PUT("/:id", h.UpdatePipe)
	pipe.GET("/all", h.GetPipes)

	pipe.GET("/:id/bookmarks", h.GetBookmarks)
	pipe.POST("/:id/bookmark", h.CreateBookmark)
	pipe.GET("/:id/bookmark/:bmId", h.GetBookmark)
	pipe.DELETE("/:id/bookmark/:bmId", h.DeleteBookmark)

}

func TestCaller(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Routing works fine",
	})
}
