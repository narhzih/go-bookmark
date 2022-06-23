package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/trencetech/mypipe-api/db"
	"gitlab.com/trencetech/mypipe-api/db/model"
	"gitlab.com/trencetech/mypipe-api/pkg/helpers"
)

func (h *Handler) EmailSignUp(c *gin.Context) {
	singUpReq := struct {
		Username    string `json:"username" binding:"required"`
		ProfileName string `json:"profile_name" binding:"required"`
		Email       string `json:"email" binding:"required"`
		Password    string `json:"password" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&singUpReq); err != nil {
		errMessage := helpers.ParseErrorMessage(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": errMessage,
			"err":     err.Error(),
		})
		return
	}

	hashedPassword, err := helpers.HashPassword(singUpReq.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went very wrong",
		})
		return
	}

	userStruct := model.User{
		Username:    singUpReq.Username,
		ProfileName: singUpReq.ProfileName,
		Email:       singUpReq.Email,
	}

	user, err := h.service.DB.CreateUserByEmail(userStruct, hashedPassword, "DEFAULT")
	if err != nil {
		if err == db.ErrRecordExists {
			h.logger.Err(err).Msg(err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Email has already been taken",
			})
			return
		}
		h.logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred while trying to register user",
			"err":     err.Error(),
		})
		return
	}

	var accountVerification model.AccountVerification
	accountVerification.UserID = user.ID
	accountVerification.Token = helpers.RandomToken(7)
	accountVerification.ExpiresAt = time.Now().Add(7200 * time.Second).Format(time.RFC3339Nano)
	accountVerification, err = h.service.DB.CreateVerification(accountVerification)
	if err != nil {
		h.logger.Err(err).Msg("An error occurred while trying to generate token details ")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"err":     err.Error(),
		})
		return
	}
	err = h.service.Mailer.SendVerificationEmail([]string{user.Email}, accountVerification.Token)
	if err != nil {
		h.logger.Err(err).Msg("An error occurred while trying to send email")
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Account created successfully. Please check your email for verification code",
		"data": map[string]interface{}{
			"v_token": accountVerification.Token,
		},
	})

}

func (h *Handler) VerifyAccount(c *gin.Context) {
	token := c.Param("token")
	if len(token) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid verification token provided",
		})
		return
	}

	tokenFromDB, err := h.service.DB.GetAccountVerificationByToken(token)
	if err != nil {
		if err == db.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Invalid verification token provided",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"err":     err.Error(),
		})
		return
	}

	parsedTime, _ := time.Parse(time.RFC3339Nano, tokenFromDB.CreatedAt)
	if tokenFromDB.Used == true || time.Now().Sub(parsedTime).Hours() > 2 {
		/** TODO: Generate a new token and send to user email and make sure that
		 * the user tha token is to be generated for is not previously verified
		 */
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Token has expired",
		})
		return
	}

	// Check if the user still exists
	user, err := h.service.DB.GetUserById(int(tokenFromDB.UserID))
	if err != nil {
		if err == db.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "User with the provided token was not found in our record",
			})
		}
	}
	user.EmailVerified = true
	user, err = h.service.MarkUserAsVerified(user, tokenFromDB.Token)
	if err != nil {
		h.logger.Err(err).Msg("Error occurred while verifying user")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"err":     err.Error(),
		})
		return
	}

	authToken, err := h.service.IssueAuthToken(user)
	if err != nil {
		h.logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred while trying to automatically log in. Please, log in manually",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Verification successful!",
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

}

func (h *Handler) EmailLogin(c *gin.Context) {
	loginReq := struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		errMessage := helpers.ParseErrorMessage(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": errMessage,
			"err":     err.Error(),
		})
		return
	}
	user, err := h.service.DB.GetUserByEmail(loginReq.Email)
	if err != nil {
		h.logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "User with specified email not found",
		})
		return
	}

	userAndAuth, err := h.service.DB.GetUserAndAuth(user)
	if err != nil {
		h.logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred while trying to sign in user",
		})
		return
	}
	verifyOk, verifyErr := helpers.VerifyPassword(loginReq.Password, userAndAuth.HashedPassword, userAndAuth.Origin)
	if verifyOk {
		authToken, err := h.service.IssueAuthToken(userAndAuth.User)
		if err != nil {
			h.logger.Err(err).Msg(err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Error occurred while trying to sign user in",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Sign in successful!",
			"data": map[string]interface{}{
				"token":         authToken.AccessToken,
				"refresh_token": authToken.RefreshToken,
				"expires_at":    authToken.ExpiresAt,
				"user": map[string]interface{}{
					"id":       user.ID,
					"username": user.Username,
					"email":    user.Email,
				},
			},
		})
	} else {
		//h.logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": verifyErr.Error(),
		})
		return
	}

	// Validate the password provided by the user

}

func (h *Handler) SignInWithGoogle(c *gin.Context) {
	var user model.User
	var isNewUser bool = false

	signInReq := struct {
		TokenString string `json:"token_string" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&signInReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	claims, err := h.service.ValidateGoogleJWT(signInReq.TokenString, c.Query("device"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid GoogleJWT",
		})
		return
	}
	h.logger.Info().Msg("google jwt validation successful")
	user, err = h.service.DB.GetUserByEmail(claims.Email)
	if err != nil {
		if err == db.ErrNoRecord {
			// Create a new user account
			h.logger.Info().Msg(fmt.Sprintf("username is %+v and email is %+v", claims.FirstName, claims.Email))
			isNewUser = true
			userCred := model.User{
				Username: claims.FirstName + " " + claims.LastName,
				Email:    claims.Email,
			}
			user, err = h.service.DB.CreateUserByEmail(userCred, "", "GOOGLE")
			if err != nil {
				h.logger.Err(err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Error occurred while trying to register user",
				})
				return
			}
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Error occurred while trying to register user",
			})
			return
		}

	}
	authToken, err := h.service.IssueAuthToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred while trying to sign user in",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Authentication successful!",
		"data": map[string]interface{}{
			"newUser":       isNewUser,
			"token":         authToken.AccessToken,
			"refresh_token": authToken.RefreshToken,
			"user": map[string]interface{}{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
			},
		},
	})

}

func (h *Handler) SignUpWithGoogle(c *gin.Context) {
	signUpReq := struct {
		TokenString string `json:"token_string" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&signUpReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}
	claims, err := h.service.ValidateGoogleJWT(signUpReq.TokenString, c.Query("device"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid GoogleJWT",
		})
		return
	}
	userCred := model.User{
		Username: claims.FirstName + " " + claims.LastName,
		Email:    claims.Email,
	}
	user, err := h.service.DB.CreateUser(userCred)
	if err != nil {
		if err == db.ErrRecordExists {
			// This means the user has already registered
			// Automatically log the user into the system
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Email has already been taken. Please provide a unique email",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred while trying to register user",
		})
		return
	}

	authToken, err := h.service.IssueAuthToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred while trying to log user in",
		})
	}
	c.JSON(http.StatusAccepted, gin.H{
		"message": "You can start organizing your life right away!",
		"data": map[string]interface{}{
			"token":         authToken.AccessToken,
			"refresh_token": authToken.RefreshToken,
			"user": map[string]interface{}{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
			},
		},
	})
}

func (h *Handler) ForgotPassword(c *gin.Context) {
	req := struct {
		Email string `json:"email" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		errMessage := helpers.ParseErrorMessage(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": errMessage,
		})
		return
	}

	user, err := h.service.DB.GetUserByEmail(req.Email)
	if err != nil {
		if err == db.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "email does not match any account in our record",
			})
			return
		}

		h.logger.Err(err).Msg("An error occurred while trying to get user")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"err":     err.Error(),
		})
		return
	}
	token := helpers.RandomToken(8)
	passwordReset, err := h.service.DB.CreatePasswordResetRecord(user, token)
	if err != nil {
		h.logger.Err(err).Msg("An error occurred while trying to send password reset token")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"err":     err.Error(),
		})
		return
	}

	err = h.service.Mailer.SendPasswordResetToken([]string{user.Email}, passwordReset.Token)
	if err != nil {
		h.logger.Err(err).Msg("An error occurred while trying to send password reset token")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"err":     err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Please check your email for instructions on how to reset your password",
	})
}

func (h *Handler) VerifyPasswordResetToken(c *gin.Context) {
	token := c.Param("token")
	if len(token) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid token provided",
		})
		return
	}

	passwordReset, err := h.service.DB.GetPasswordResetRecord(token)
	if err != nil {
		if err == db.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Invalid token provided",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"err":     err.Error(),
		})
		return
	}

	user, err := h.service.DB.GetUserById(int(passwordReset.UserID))
	if err != nil {
		if err == db.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "User with the attached token does not exist",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"err":     err.Error(),
		})
		return
	}

	parsedTime, _ := time.Parse(time.RFC3339Nano, passwordReset.CreatedAt)
	if time.Now().Sub(parsedTime).Hours() > 2 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Token has expired",
		})
		return
	}

	passwordReset, err = h.service.DB.UpdatePasswordResetRecord(passwordReset.Token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"err":     err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Proceed to set your new password",
		"data": map[string]interface{}{
			"user": map[string]interface{}{
				"id":       user.ID,
				"email":    user.Email,
				"username": user.Username,
			},
			"token": passwordReset.Token,
		},
	})

}

func (h *Handler) ResetPassword(c *gin.Context) {
	token := c.Param("token")
	if len(token) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid token provided",
		})
		return
	}

	// Check if token exists in the DB
	passwordReset, err := h.service.DB.GetPasswordResetRecord(token)
	if err != nil {
		if err == db.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Invalid token provided",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"err":     err.Error(),
		})
		return
	}

	if passwordReset.Validated != true {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid token provided because it's not validated",
		})
		return
	}

	// Check if request made is valid
	resetReq := struct {
		Password string `json:"password" binding:"required"`
	}{}

	if err = c.ShouldBindJSON(&resetReq); err != nil {
		errMessage := helpers.ParseErrorMessage(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": errMessage,
		})
	}

	// check if user with provided email is found
	user, err := h.service.DB.GetUserById(int(passwordReset.UserID))
	if err != nil {
		if err == db.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "user with attached token not found",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"err":     err.Error(),
		})
		return
	}

	// Update the user's password
	hashedPassword, err := helpers.HashPassword(resetReq.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something" +
				"" +
				" went very wrong",
		})
		return
	}
	err = h.service.DB.UpdateUserPassword(int(user.ID), hashedPassword)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went very wrong",
			"err":     err.Error(),
		})
		return
	}

	// Password update successfully, delete the record and generate new login token
	err = h.service.DB.DeletePasswordResetRecord(passwordReset.Token)
	if err != nil {
		h.logger.Err(err).Msg("An error occurred while trying to delete password reset record")
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully. Please proceed to login with your new password",
	})

	//authToken, err := h.service.IssueAuthToken(user)
	//if err != nil {
	//	h.logger.Err(err).Msg(err.Error())
	//	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
	//		"message": "Error occurred while trying to sign user in. Please try to login manually",
	//	})
	//	return
	//}
	//
	//c.JSON(http.StatusOK, gin.H{
	//	"message": "Password updated successfully",
	//	"data": map[string]interface{}{
	//		"token":         authToken.AccessToken,
	//		"refresh_token": authToken.RefreshToken,
	//		"expires_at":    authToken.ExpiresAt,
	//		"user": map[string]interface{}{
	//			"id":           user.ID,
	//			"email":        user.Email,
	//			"profile_name": user.ProfileName,
	//			"username":     user.Username,
	//			"cover_photo":  user.CovertPhoto,
	//		},
	//	},
	//})

}
