package api

import (
	"github.com/rs/zerolog"
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
	// updateReq := struct {
	// 	Username    string `json:"username"`
	// 	CoverPhoto  string `json:"cover_photo"`
	// 	ProfileName string `json:"profile_name"`
	// }{}

	// // Check if an image was uploaded
	// file, header, err := c.Request.FormFile("cover_photo")
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusBadGateway, gin.H {
	// 		"message": fmt.Sprintf("file err : %s", err.Error()),
	// 	})
	// } else {
	// 	h.logger.Info().Msg("An image upload was actually detected")
	// }

	//if err := c.ShouldBindJSON(&updateReq); err != nil {
	//	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
	//		"message": "Invalid request body",
	//	})
	//	return
	//}
	//
	//updatedUser := model.User{
	//	ID:          c.GetInt64(KeyUserId),
	//	Username:    updateReq.Username,
	//	CovertPhoto: updateReq.CoverPhoto,
	//	ProfileName: updateReq.ProfileName,
	//}
	//user, err := h.service.DB.UpdateUser(updatedUser)
	//if err != nil {
	//	h.logger.Err(err).Msg(err.Error())
	//	if err == db.ErrNoRecord {
	//		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
	//			"message": "Could not  update user because user was not found",
	//		})
	//		return
	//	}
	//
	//	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
	//		"message": "An error occurred when trying to update user",
	//	})
	//	return
	//}
	//
	//c.JSON(http.StatusCreated, gin.H{
	//	"message": "Profile updated successfully",
	//	"data": map[string]interface{}{
	//		"user": map[string]interface{}{
	//			"id":          user.ID,
	//			"username":    user.Username,
	//			"cover_photo": user.CovertPhoto,
	//		},
	//	},
	//})
}

func (h *Handler) UploadCoverPhoto(c *gin.Context) {
	logger := zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()
	uploadInformation := service.FileUploadInformation{
		Logger:        logger,
		Ctx:           c,
		FileInputName: "cover_photo",
		Type: "user",
	}
	photoUrl, err := service.UploadToCloudinary(uploadInformation)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H {
			"message": "An error occurred when trying to save user image",
		})
	}

	logger.Info().Msg(photoUrl)
	c.JSON(http.StatusOK, gin.H{
		"message": "Image uploaded successfully "+photoUrl,
	})
}
