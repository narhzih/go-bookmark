package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/gin-gonic/gin"
	"github.com/mypipeapp/mypipeapi/cmd/api/helpers"
	"github.com/mypipeapp/mypipeapi/cmd/api/internal"
	"github.com/mypipeapp/mypipeapi/cmd/api/middlewares"
	"github.com/mypipeapp/mypipeapi/cmd/api/services"
	"github.com/mypipeapp/mypipeapi/db/actions/postgres"
	"github.com/mypipeapp/mypipeapi/db/models"
	"net/http"
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
	req := struct {
		Username    string `form:"username" json:"username,omitempty"`
		Email       string `form:"email" json:"email,omitempty"`
		ProfileName string `form:"profile_name" json:"profile_name,omitempty"`
		TwitterId   string `form:"twitter_id" json:"twitter_id,omitempty"`
	}{}

	if err := c.ShouldBind(&req); err != nil {
		errMessage := helpers.ParseErrorMessage(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": errMessage,
			"err":     err.Error(),
		})
		return
	}

	// fetch the current logged in user first
	user, _ := h.app.Repositories.User.GetUserById(c.GetInt64(middlewares.KeyUserId))
	userBytes, err := json.Marshal(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "our server encountered an error",
			"err":     "Could not bind existing user body",
		})
		return
	}

	reqBytes, err := json.Marshal(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "our server encountered an error",
			"err":     "Could not bind request body",
		})
		return
	}

	patchBody, err := jsonpatch.MergePatch(userBytes, reqBytes)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "our server encountered an error",
			"err":     "Could not merge user and request bytes",
		})
		return
	}

	err = json.Unmarshal(patchBody, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "our server encountered an error",
			"err":     "Could not unmarhshal user body",
		})
		return
	}

	// if there's a cover photo sent alongside the request, upload it and update the
	// request body
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

		user.CovertPhoto = photoUrl
	}
	user, err = h.app.Repositories.User.UpdateUser(user)

	if err != nil {
		switch {
		case errors.Is(err, postgres.ErrNoRecord):
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Could not update user because user was not found",
				"err":     err.Error(),
			})
			return
		case errors.Is(err, postgres.ErrDuplicateEmail):
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "user with email already exits",
				"err":     err.Error(),
			})
			return
		case errors.Is(err, postgres.ErrDuplicateUsername):
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "user with username already exits",
				"err":     err.Error(),
			})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "An error occurred while trying to update user",
				"err":     err.Error(),
			})
			return
		}

	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Profile updated successfully",
		"data": map[string]interface{}{
			"user": user,
		},
	})
}

func (h userHandler) UploadCoverPhoto(c *gin.Context) {
	var user models.User
	uploadInformation := services.FileUploadInformation{
		Logger:        h.app.Logger,
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
	user, _ = h.app.Repositories.User.GetUserById(c.GetInt64(middlewares.KeyUserId))
	user.CovertPhoto = photoUrl
	user, err = h.app.Repositories.User.UpdateUser(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred while trying update user profile",
			"err":     err.Error(),
		})
		return
	}
	h.app.Logger.Info().Msg(photoUrl)
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
	user, err := h.app.Repositories.User.GetUserById(userID)
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
	userAndAuth, err := h.app.Repositories.User.GetUserAndAuth(c.GetInt64(middlewares.KeyUserId))
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
		err = h.app.Repositories.User.UpdateUserPassword(userAndAuth.User.ID, newPasswordHash)
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
