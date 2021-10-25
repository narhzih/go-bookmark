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
}

func TestCaller(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Routing works fine",
	})
}
