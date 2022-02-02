package api

import (
	"fmt"
	"github.com/rs/zerolog"
	"gitlab.com/trencetech/mypipe-api/db"
	"gitlab.com/trencetech/mypipe-api/pkg/helpers"
	"gitlab.com/trencetech/mypipe-api/pkg/service"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gitlab.com/trencetech/mypipe-api/db/model"
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

	updatedUser := model.User{
		ID: c.GetInt64(KeyUserId),
	}
	username := c.PostForm("username")
	profileName := c.PostForm("profile_name")

	// Check if there's already a user with the same username
	if len(username) > 0 {
		userWithUsername, err := h.service.DB.GetUserByUsername(username)
		if err != nil {
			if err != db.ErrNoRecord {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "An error occurred while getting user",
					"err":     err.Error(),
				})
				return
			}
		}
		if err != db.ErrNoRecord && userWithUsername.ID != c.GetInt64(KeyUserId) {
			// This means there's another user that is not
			// the user making the same request who has the same username
			h.logger.Info().Msg(userWithUsername.Email)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "username has already been taken by another user",
			})
			return
		}

		updatedUser.Username = username
	}
	if len(profileName) > 0 {
		updatedUser.ProfileName = profileName
	}

	// At this point, username and profile_name validation passes
	// Try uploading the image to Cloud if any was parsed
	file, _, err := c.Request.FormFile("cover_photo")
	if err != nil {
		if err != http.ErrMissingFile {
			h.logger.Err(err).Msg(fmt.Sprintf("file err : %s", err.Error()))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "An error occurred from here",
				"err":     err.Error(),
			})
			return
		}

	}
	if file != nil {
		// This means a file was uploaded with the request
		// Try uploading it to Cloudinary
		uploadInformation := service.FileUploadInformation{
			Logger:        h.logger,
			Ctx:           c,
			FileInputName: "cover_photo",
			Type:          "user",
		}
		photoUrl, err := service.UploadToCloudinary(uploadInformation)
		h.logger.Info().Msg("The cover photo is " + photoUrl)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "An error occurred when trying to save user image",
				"err":     err.Error(),
			})
			return
		}

		updatedUser.CovertPhoto = photoUrl
	}
	h.logger.Info().Msg("Photo url after property setting is - " + updatedUser.CovertPhoto)
	user, err := h.service.DB.UpdateUser(updatedUser)
	h.logger.Info().Msg("Photo url after normal upload is - " + user.CovertPhoto)

	if err != nil {
		h.logger.Err(err).Msg(err.Error())
		if err == db.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Could not  update user because user was not found",
				"err":     err.Error(),
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred when trying to update user",
			"err":     err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Profile updated successfully",
		"data": map[string]interface{}{
			"user": user,
		},
	})
}

func (h *Handler) UploadCoverPhoto(c *gin.Context) {
	logger := zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()
	uploadInformation := service.FileUploadInformation{
		Logger:        logger,
		Ctx:           c,
		FileInputName: "cover_photo",
		Type:          "user",
	}
	photoUrl, err := service.UploadToCloudinary(uploadInformation)
	if err != nil {
		if err == http.ErrMissingFile {
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"message": "No file was uploaded. Please select a file to upload",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred when trying to save user image",
		})
		return
	}
	updatedUserModel := model.User{
		ID:          c.GetInt64(KeyUserId),
		CovertPhoto: photoUrl,
	}
	user, err := h.service.DB.UpdateUser(updatedUserModel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred while trying update user profile",
			"err":     err.Error(),
		})
		return
	}
	logger.Info().Msg(photoUrl)
	c.JSON(http.StatusOK, gin.H{
		"message": "Image uploaded successfully",
		"data": map[string]interface{}{
			"user": map[string]interface{}{
				"id":           user.ID,
				"cover_photo":  user.CovertPhoto,
				"email":        user.Email,
				"profile_name": user.ProfileName,
				"username":     user.Username,
			},
		},
	})
}

func (h *Handler) ChangePassword(c *gin.Context) {
	reqBody := struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		errMessage := helpers.ParseErrorMessage(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": errMessage,
			"error":   err.Error(),
		})
		return
	}
	if reqBody.Password != reqBody.ConfirmPassword {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Passwords do not match",
		})
		return
	}
	userID := c.GetInt64(KeyUserId)
	user, err := h.service.DB.GetUserById(int(userID))
	if err != nil {
		h.logger.Err(err).Msg(err.Error())
		if err == db.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "You have to log in to be able to perform this operation",
			})
			return
		}
		h.logger.Err(err).Msg("An error occurred while trying to get user from the database")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"error":   err.Error(),
		})
		return
	}
	userAndAuth, err := h.service.DB.GetUserAndAuth(user)
	if err != nil {
		h.logger.Err(err).Msg("An error occurred while trying to get user from the database")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"error":   err.Error(),
		})
		return
	}

	if helpers.VerifyPassword(reqBody.CurrentPassword, userAndAuth.HashedPassword) {
		newPasswordHash, err := helpers.HashPassword(reqBody.Password)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong",
				"error":   err.Error(),
			})
			return
		}
		err = h.service.DB.UpdateUserPassword(int(userAndAuth.User.ID), newPasswordHash)
		if err != nil {
			h.logger.Err(err).Msg(err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong",
				"error":   err.Error(),
			})
			return
		}
		authToken, err := h.service.IssueAuthToken(user)
		if err != nil {
			h.logger.Err(err).Msg(err.Error())
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"message": "Password changed successfully",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Password changed successfully",
			"data": map[string]interface{}{
				"token":         authToken.AccessToken,
				"refresh_token": authToken.RefreshToken,
				"expires_at":    authToken.ExpiresAt,
				"user": map[string]interface{}{
					"id":           user.ID,
					"username":     user.Username,
					"email":        user.Email,
					"profile name": user.ProfileName,
					"cover_photo":  user.CovertPhoto,
				},
			},
		})

	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Current password incorrect",
		})
	}
}
