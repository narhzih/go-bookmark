package api

import (
	"fmt"
	"github.com/rs/zerolog"
	"gitlab.com/trencetech/mypipe-api/db"
	"gitlab.com/trencetech/mypipe-api/pkg/service"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gitlab.com/trencetech/mypipe-api/db/model"
)

func (h *Handler) OnboardUser(c *gin.Context) {
	// TODO : Implement route to onboard user
	// i.e Set username and twitter_handle
	// var err error

	// onboardRequest := struct {
	// 	Username    string `json:"username" binding:"required"`
	// 	ProfileName string `json:"twitter_handle" binding:"required"`
	// 	Email       string `json:"email" binding:"required"`
	// 	Password    string `json:"password" binding:"required"`
	// }{}

	// if err := c.ShouldBindJSON(&onboardRequest); err != nil {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
	// 		"message": "Invalid request body",
	// 	})
	// 	return
	// }

	// newUserReq := model.User{
	// 	Email:       onboardRequest.Email,
	// 	Username:    onboardRequest.Username,
	// 	ProfileName: onboardRequest.ProfileName,
	// }

	// _, err = h.service.DB.CreateUserByEmail(newUserReq, onboardRequest.Password)
	// if err != nil {
	// 	// Other error checks will be implemented soon
	// 	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
	// 		"message": "An error occurred while trying to create user account",
	// 	})
	// }

	/**
	*	TODO:
	*	Immediately user account is created, send a verification email to the user account
	 */

	c.JSON(http.StatusCreated, gin.H{
		"message": "Account created successfully. Please verify your account by following the instructions sent to your email",
	})
}

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
					"message": "An error occurred",
					"err":     err.Error(),
				})
				return
			}
		}
		if err != db.ErrNoRecord && userWithUsername.ID != c.GetInt64(KeyUserId) {
			// This means there's another user that is not
			// the user making the same request who has the same username
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
				"message": "An error occurred",
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
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "An error occurred when trying to save user image",
				"err":     err.Error(),
			})
			return
		}

		updatedUser.CovertPhoto = photoUrl
	}
	user, err := h.service.DB.UpdateUser(updatedUser)
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
			"user": map[string]interface{}{
				"id":          user.ID,
				"username":    user.Username,
				"cover_photo": user.CovertPhoto,
			},
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
