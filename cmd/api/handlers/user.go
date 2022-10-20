package handlers

import (
	"fmt"
	"github.com/mypipeapp/mypipeapi/cmd/api/helpers"
	"github.com/mypipeapp/mypipeapi/cmd/api/internal"
	"github.com/mypipeapp/mypipeapi/cmd/api/middlewares"
	"github.com/mypipeapp/mypipeapi/cmd/api/services"
	"github.com/mypipeapp/mypipeapi/db/actions/postgres"
	"github.com/mypipeapp/mypipeapi/db/models"
	"github.com/rs/zerolog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	UserProfile(c *gin.Context)
	EditProfile(c *gin.Context)
	UploadCoverPhoto(c *gin.Context)
	ChangePassword(c *gin.Context)
}

type userHandler struct {
	app internal.Application
}

func NewUserHandler(app internal.Application) UserHandler {
	return userHandler{app: app}
}

func (h userHandler) UserProfile(c *gin.Context) {
	var userProfile models.Profile
	var err error
	userID := c.GetInt64(middlewares.KeyUserId)
	userProfile, err = h.app.Services.GetUserProfileInformation(userID)

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

func (h userHandler) EditProfile(c *gin.Context) {

	updatedUser := models.User{
		ID: c.GetInt64(middlewares.KeyUserId),
	}
	username := c.PostForm("username")
	profileName := c.PostForm("profile_name")

	// Check if there's already a user with the same username
	if len(username) > 0 {
		userWithUsername, err := h.app.Repositories.User.GetUserByUsername(username)
		if err != nil {
			if err != postgres.ErrNoRecord {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "An error occurred while getting user",
					"err":     err.Error(),
				})
				return
			}
		}
		if err != postgres.ErrNoRecord && userWithUsername.ID != c.GetInt64(middlewares.KeyUserId) {
			// This means there's another user that is not
			// the user making the same request who has the same username
			h.app.Logger.Info().Msg(userWithUsername.Email)
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
			h.app.Logger.Err(err).Msg(fmt.Sprintf("file err : %s", err.Error()))
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
		uploadInformation := services.FileUploadInformation{
			Logger:        h.app.Logger,
			Ctx:           c,
			FileInputName: "cover_photo",
			Type:          "user",
		}
		photoUrl, err := services.UploadToCloudinary(uploadInformation)
		h.app.Logger.Info().Msg("The cover photo is " + photoUrl)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "An error occurred when trying to save user image",
				"err":     err.Error(),
			})
			return
		}

		updatedUser.CovertPhoto = photoUrl
	}
	h.app.Logger.Info().Msg("Photo url after property setting is - " + updatedUser.CovertPhoto)
	user, err := h.app.Repositories.User.UpdateUser(updatedUser)
	h.app.Logger.Info().Msg("Photo url after normal upload is - " + user.CovertPhoto)

	if err != nil {
		h.app.Logger.Err(err).Msg(err.Error())
		if err == postgres.ErrNoRecord {
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

func (h userHandler) UploadCoverPhoto(c *gin.Context) {
	logger := zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()
	uploadInformation := services.FileUploadInformation{
		Logger:        logger,
		Ctx:           c,
		FileInputName: "cover_photo",
		Type:          "user",
	}
	photoUrl, err := services.UploadToCloudinary(uploadInformation)
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
	updatedUserModel := models.User{
		ID:          c.GetInt64(middlewares.KeyUserId),
		CovertPhoto: photoUrl,
	}
	user, err := h.app.Repositories.User.UpdateUser(updatedUserModel)
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

func (h userHandler) ChangePassword(c *gin.Context) {
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
	userID := c.GetInt64(middlewares.KeyUserId)
	user, err := h.app.Repositories.User.GetUserById(int(userID))
	if err != nil {
		h.app.Logger.Err(err).Msg(err.Error())
		if err == postgres.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "You have to log in to be able to perform this operation",
			})
			return
		}
		h.app.Logger.Err(err).Msg("An error occurred while trying to get user from the database")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"error":   err.Error(),
		})
		return
	}
	userAndAuth, err := h.app.Repositories.User.GetUserAndAuth(user)
	if err != nil {
		h.app.Logger.Err(err).Msg("An error occurred while trying to get user from the database")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"error":   err.Error(),
		})
		return
	}
	verifyOk, verifyErr := helpers.VerifyPassword(reqBody.CurrentPassword, userAndAuth.HashedPassword, userAndAuth.Origin)
	if verifyOk {
		newPasswordHash, err := helpers.HashPassword(reqBody.Password)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong",
				"error":   err.Error(),
			})
			return
		}
		err = h.app.Repositories.User.UpdateUserPassword(int(userAndAuth.User.ID), newPasswordHash)
		if err != nil {
			h.app.Logger.Err(err).Msg(err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong",
				"error":   err.Error(),
			})
			return
		}
		authToken, err := h.app.Services.IssueAuthToken(user)
		if err != nil {
			h.app.Logger.Err(err).Msg(err.Error())
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
			"message": verifyErr.Error(),
		})
	}
}
