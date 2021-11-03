package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gowagr/mypipe-api/db/model"
)

func (h *Handler) UserProfile(c *gin.Context) {
	var userProfile model.Profile
	var err error
	userID := c.GetInt64(KeyUserId)
	userProfile, err = h.service.GetUserProfileInformation(userID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "An error occurred while trying to get user profile information",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile fetched successfully",
		"data": map[string]interface{}{
			"user":      userProfile.User,
			"pipes":     userProfile.Pipes,
			"bookmarks": userProfile.Bookmarks,
		},
	})
}
