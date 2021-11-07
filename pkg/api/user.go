package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gowagr/mypipe-api/db"
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

func (h *Handler) EditProfile(c *gin.Context) {
	updateReq := struct {
		Username   string `json:"username"`
		CoverPhoto string `json:"cover_photo"`
	}{}

	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	updatedUser := model.User{
		ID:          c.GetInt64(KeyUserId),
		Username:    updateReq.Username,
		CovertPhoto: updateReq.CoverPhoto,
	}
	user, err := h.service.DB.UpdateUser(updatedUser)
	if err != nil {
		if err == db.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Could not  update user because user was not found",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred when trying to update user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Profile updated successfully",
		"data": map[string]interface{}{
			"user": map[string]interface{}{
				"id":          user.ID,
				"username":    user.Username,
				"cover_photo": user.CovertPhoto,
			},
		},
	})
}
