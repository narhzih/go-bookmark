package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gitlab.com/trencetech/mypipe-api/pkg/service"
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

	// This particular route is just for testing purposes
	// it will be removed later when the front end starts
	// sending in actual GoogleJWT
	routeGroup.POST("/sign-up", h.EmailSignUp)
	routeGroup.POST("/sign-in", h.EmailLogin)
	routeGroup.POST("/verify-account/:token", h.VerifyAccount)
	routeGroup.POST("/forgot-password", h.ForgotPassword)
	routeGroup.POST("/verify-reset-token/:token", h.VerifyPasswordResetToken)
	routeGroup.POST("/reset-password/:token", h.ResetPassword)

	routeGroup.POST("/google-auth", h.SignInWithGoogle)
	// routeGroup.POST("/google/sign-in", h.SignInWithGoogle)

	authApi := routeGroup.Group("/auth")
	authApi.Use(AuthRequired(h.service.JWTConfig.Key, h.logger))
	authApi.GET("/auth-test", TestCaller)
	authApi.POST("/twitter/connect-account", h.ConnectTwitterAccount)

	pipe := routeGroup.Group("/pipe")
	pipe.Use(AuthRequired(h.service.JWTConfig.Key, h.logger))
	pipe.POST("/", h.CreatePipe)
	pipe.GET("/:id", h.GetPipe)
	pipe.POST("/:id/share", h.SharePipe)
	pipe.PUT("/:id", h.UpdatePipe)
	pipe.DELETE("/:id", h.DeletePipe)
	pipe.GET("/all", h.GetPipes)
	pipe.GET("/all/steroids", h.GetPipeWithResource)
	pipe.GET("/preview", h.PreviewPipe)
	pipe.POST("/add-pipe", h.AddPipe)

	pipe.GET("/:id/bookmarks", h.GetBookmarks)
	pipe.POST("/:id/bookmark", h.CreateBookmark)
	pipe.GET("/:id/bookmark/:bmId", h.GetBookmark)
	pipe.DELETE("/:id/bookmark/:bmId", h.DeleteBookmark)

	user := routeGroup.Group("/user")
	user.Use(AuthRequired(h.service.JWTConfig.Key, h.logger))
	user.GET("/profile", h.UserProfile)
	user.PATCH("/profile", h.EditProfile)
	user.PATCH("/profile/change-password", h.ChangePassword)
	user.POST("/profile/cover-photo", h.UploadCoverPhoto)

	notification := routeGroup.Group("/notifications")
	notification.Use(AuthRequired(h.service.JWTConfig.Key, h.logger))
	notification.GET("/", h.GetNotifications)
	notification.POST("/update-device-tokens", h.UpdateUserDeviceTokens)
	notification.GET("/:notificationId", h.GetNotification)

	parser := routeGroup.Group("/parse-link")
	//parser.Use(AuthRequired(h.service.JWTConfig.Key, h.logger))
	parser.POST("/twitter", h.TwitterLinkParser)
	parser.POST("/youtube", h.YoutubeLinkParser)
	parser.POST("/others", h.ParseLink)

	bot := routeGroup.Group("/bot")
	bot.POST("/twitter/add-to-pipe", h.BotAddToPipe)

}

func TestCaller(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Routing works fine",
	})
}
