package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mypipeapp/mypipeapi/cmd/api/helpers"
	"github.com/mypipeapp/mypipeapi/cmd/api/internal"
	"github.com/mypipeapp/mypipeapi/cmd/api/middlewares"
	"github.com/mypipeapp/mypipeapi/cmd/api/models/response"
	"github.com/mypipeapp/mypipeapi/db/actions/postgres"
	"github.com/mypipeapp/mypipeapi/db/models"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	EmailSignUp(c *gin.Context)
	EmailLogin(c *gin.Context)
	VerifyAccount(c *gin.Context)
	ForgotPassword(c *gin.Context)
	VerifyPasswordResetToken(c *gin.Context)
	ResetPassword(c *gin.Context)
	SignInWithGoogle(c *gin.Context)
	ConnectTwitterAccount(c *gin.Context)
	GetConnectedTwitterAccount(c *gin.Context)
	DisconnectTwitterAccount(c *gin.Context)
}

type authHandler struct {
	app internal.Application
}

func NewAuthHandler(app internal.Application) AuthHandler {
	return authHandler{
		app: app,
	}
}

func (h authHandler) EmailSignUp(c *gin.Context) {
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

	userStruct := models.User{
		Username:    singUpReq.Username,
		ProfileName: singUpReq.ProfileName,
		Email:       singUpReq.Email,
	}

	user, err := h.app.Repositories.User.CreateUserByEmail(userStruct, hashedPassword, "DEFAULT")
	if err != nil {
		if err == postgres.ErrRecordExists {
			h.app.Logger.Err(err).Msg(err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Email has already been taken",
			})
			return
		}
		h.app.Logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred while trying to register user",
			"err":     err.Error(),
		})
		return
	}

	var accountVerification models.AccountVerification
	accountVerification.UserID = user.ID
	accountVerification.Token = helpers.RandomToken(7)
	accountVerification.ExpiresAt = time.Now().Add(7200 * time.Second)
	accountVerification, err = h.app.Repositories.AccountVerification.CreateVerification(accountVerification)
	if err != nil {
		h.app.Logger.Err(err).Msg("An error occurred while trying to generate token details ")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"err":     err.Error(),
		})
		return
	}
	err = h.app.Services.Mailer.SendVerificationEmail([]string{user.Email}, accountVerification.Token)
	if err != nil {
		h.app.Logger.Err(err).Msg("An error occurred while trying to send email")
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Account created successfully. Please check your email for verification code",
		"data": map[string]interface{}{
			"v_token": accountVerification.Token,
		},
	})

}

func (h authHandler) VerifyAccount(c *gin.Context) {
	token := c.Param("token")
	if len(token) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid verification token provided",
		})
		return
	}

	tokenFromDB, err := h.app.Repositories.AccountVerification.GetAccountVerificationByToken(token)
	if err != nil {
		if err == postgres.ErrNoRecord {
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

	//parsedTime, _ := time.Parse(time.RFC3339Nano, tokenFromDB.CreatedAt)
	if tokenFromDB.Used == true || time.Now().Sub(tokenFromDB.CreatedAt).Hours() > 2 {
		/** TODO: Generate a new token and send to user email and make sure that
		 * the user tha token is to be generated for is not previously verified
		 */
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Token has expired",
		})
		return
	}

	// Check if the user still exists
	user, err := h.app.Repositories.User.GetUserById(tokenFromDB.UserID)
	if err != nil {
		if err == postgres.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "User with the provided token was not found in our record",
			})
		}
	}
	user.EmailVerified = true
	user, err = h.app.Services.MarkUserAsVerified(user, tokenFromDB.Token)
	if err != nil {
		h.app.Logger.Err(err).Msg("Error occurred while verifying user")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"err":     err.Error(),
		})
		return
	}

	authToken, err := h.app.Services.IssueAuthToken(user)
	if err != nil {
		h.app.Logger.Err(err).Msg(err.Error())
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

func (h authHandler) EmailLogin(c *gin.Context) {
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
	user, err := h.app.Repositories.User.GetUserByEmail(loginReq.Email)
	if err != nil {
		h.app.Logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "User with specified email not found",
			"err":     err.Error(),
		})
		return
	}

	userAndAuth, err := h.app.Repositories.User.GetUserAndAuth(user.ID)
	if err != nil {
		h.app.Logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred while trying to sign in user",
			"err":     err.Error(),
		})
		return
	}
	verifyOk, verifyErr := helpers.VerifyPassword(loginReq.Password, userAndAuth.HashedPassword, userAndAuth.Origin)
	if verifyOk {
		authToken, err := h.app.Services.IssueAuthToken(userAndAuth.User)
		if err != nil {
			h.app.Logger.Err(err).Msg(err.Error())
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": verifyErr.Error(),
		})
		return
	}

	// Validate the password provided by the user

}

func (h authHandler) SignInWithGoogle(c *gin.Context) {
	var user models.User
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

	claims, err := h.app.Services.ValidateGoogleJWT(signInReq.TokenString, c.Query("device"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid GoogleJWT",
			"err":     err.Error(),
		})
		return
	}
	h.app.Logger.Info().Msg("google jwt validation successful")
	user, err = h.app.Repositories.User.GetUserByEmail(claims.Email)
	if err != nil {
		if err == postgres.ErrNoRecord {
			// Create a new user account
			h.app.Logger.Info().Msg(fmt.Sprintf("username is %+v and email is %+v", claims.GivenName, claims.Email))
			isNewUser = true
			userCred := models.User{
				Username:    strings.TrimSpace(claims.GivenName),
				Email:       claims.Email,
				ProfileName: strings.TrimSpace(claims.GivenName) + " " + strings.TrimSpace(claims.FamilyName),
			}
			user, err = h.app.Repositories.User.CreateUserByEmail(userCred, "", "GOOGLE")
			if err != nil {
				h.app.Logger.Err(err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Error occurred while trying to register user",
					"err":     err.Error(),
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
	authToken, err := h.app.Services.IssueAuthToken(user)
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
			"user":          user,
		},
	})

}

func (h authHandler) ForgotPassword(c *gin.Context) {
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

	user, err := h.app.Repositories.User.GetUserByEmail(req.Email)
	if err != nil {
		if err == postgres.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "email does not match any account in our record",
			})
			return
		}

		h.app.Logger.Err(err).Msg("An error occurred while trying to get user")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"err":     err.Error(),
		})
		return
	}
	token := helpers.RandomToken(8)
	passwordReset, err := h.app.Repositories.PasswordReset.CreatePasswordResetRecord(user, token)
	if err != nil {
		h.app.Logger.Err(err).Msg("An error occurred while trying to send password reset token")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"err":     err.Error(),
		})
		return
	}

	err = h.app.Services.Mailer.SendPasswordResetToken([]string{user.Email}, passwordReset.Token)
	if err != nil {
		h.app.Logger.Err(err).Msg("An error occurred while trying to send password reset token")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "our server encountered an error",
			"err":     err.Error(),
		})
		return
	}

	// for testing/development purposes, return the token as part of the response
	if os.Getenv("APP_ENV") != "prod" {
		c.JSON(http.StatusOK, gin.H{
			"message": "Please check your email for instructions on how to reset your password",
			"token":   token,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Please check your email for instructions on how to reset your password",
		})
	}
}

func (h authHandler) VerifyPasswordResetToken(c *gin.Context) {
	token := c.Param("token")
	if len(token) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid token provided",
		})
		return
	}

	passwordReset, err := h.app.Repositories.PasswordReset.GetPasswordResetRecord(token)
	if err != nil {
		if err == postgres.ErrNoRecord {
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

	user, err := h.app.Repositories.User.GetUserById(passwordReset.UserID)
	if err != nil {
		if err == postgres.ErrNoRecord {
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

	passwordReset, err = h.app.Repositories.PasswordReset.UpdatePasswordResetRecord(passwordReset.Token)
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

func (h authHandler) ResetPassword(c *gin.Context) {
	token := c.Param("token")
	if len(token) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid token provided",
		})
		return
	}

	// Check if token exists in the DB
	passwordReset, err := h.app.Repositories.PasswordReset.GetPasswordResetRecord(token)
	if err != nil {
		if err == postgres.ErrNoRecord {
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
			"message": "Invalid token provided",
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
	user, err := h.app.Repositories.User.GetUserById(passwordReset.UserID)
	if err != nil {
		if err == postgres.ErrNoRecord {
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
	err = h.app.Repositories.User.UpdateUserPassword(user.ID, hashedPassword)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went very wrong",
			"err":     err.Error(),
		})
		return
	}

	// Password update successfully, delete the record and generate new login token
	err = h.app.Repositories.PasswordReset.DeletePasswordResetRecord(passwordReset.Token)
	if err != nil {
		h.app.Logger.Err(err).Msg("An error occurred while trying to delete password reset record")
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully. Please proceed to login with your new password",
	})
}

func (h authHandler) ConnectTwitterAccount(c *gin.Context) {
	req := struct {
		AccessToken string `json:"accessToken" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&req); err != nil {
		errMessage := helpers.ParseErrorMessage(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": errMessage,
			"err":     err.Error(),
		})
		return
	}

	//h.app.Logger.Info().Msg(req.AccessToken)
	user, err := h.app.Repositories.User.GetUserById(c.GetInt64(middlewares.KeyUserId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"err":     err.Error(),
		})
	}

	//h.app.Logger.Info().Msg(fmt.Sprintf("access token is %v", req.AccessToken))
	// Trying to use the auth code to get a valid access token

	// Make request to api for user information
	twitterHttp, err := http.NewRequest(http.MethodGet, "https://api.twitter.com/2/users/me", nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "we can't proceed to connect your twitter account at the moment. Please try again soon",
		})
		return
	}
	twitterHttp.Header.Add("Authorization", "Bearer "+req.AccessToken)
	twitterResponse, err := http.DefaultClient.Do(twitterHttp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred while connecting to twitter api!",
			"err":     err.Error(),
		})
		return
	}

	if twitterResponse.StatusCode != http.StatusOK {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid access token provided",
		})
		return
	}
	var twitterUserResponse response.TwitterUserResponse
	respBody, err := io.ReadAll(twitterResponse.Body)
	json.Unmarshal(respBody, &twitterUserResponse)
	existingAcc, err := h.app.Repositories.User.GetUserByTwitterID(twitterUserResponse.Data.Id)
	if err != nil {
		switch {
		case errors.Is(err, postgres.ErrNoRecord):
			user, err = h.app.Repositories.User.ConnectToTwitter(user, twitterUserResponse.Data.Id)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Server error",
					"err":     err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Twitter account connected successfully",
				"data": map[string]interface{}{
					"user": map[string]interface{}{
						"id":           user.ID,
						"username":     user.Username,
						"email":        user.Email,
						"profile name": user.ProfileName,
						"cover_photo":  user.CovertPhoto,
						"twitter_id":   user.TwitterId,
					},
				},
			})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}
	}
	// Format error
	if existingAcc.ID != user.ID {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "This twitter account has already been connected to another account on our database.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "You have already connected your twitter account to your mypipe account. You don't have to connect again",
	})

}

func (h authHandler) GetConnectedTwitterAccount(c *gin.Context) {
	authenticatedUser, err := h.app.Repositories.User.GetUserById(c.GetInt64(middlewares.KeyUserId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}
	if authenticatedUser.TwitterId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No twitter handle connected to this account",
		})
		return
	}

	// Fetch user information from Twitter
	url := fmt.Sprintf("https://api.twitter.com/2/users/%v?user.fields=profile_image_url", authenticatedUser.TwitterId)
	twitterHttp, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "an error occurred while connecting to twitter api",
			"err":     err.Error(),
		})
		return
	}
	twitterHttp.Header.Add("Authorization", fmt.Sprintf("Bearer %v", os.Getenv("BEARER_TOKEN")))

	twitterResponse, err := http.DefaultClient.Do(twitterHttp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred while connecting to twitter api!",
			"err":     err.Error(),
		})
		return
	}
	if twitterResponse.StatusCode != http.StatusOK {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid access token provided",
		})
		return
	}
	var twitterUserResponse response.TwitterUserResponse
	respBody, err := io.ReadAll(twitterResponse.Body)
	json.Unmarshal(respBody, &twitterUserResponse)
	c.JSON(http.StatusOK, gin.H{
		"loggedInUser": authenticatedUser,
		"twitterAccount": map[string]interface{}{
			"details": map[string]interface{}{
				"username":      twitterUserResponse.Data.Username,
				"name":          twitterUserResponse.Data.Name,
				"id":            twitterUserResponse.Data.Id,
				"profile_photo": twitterUserResponse.Data.ProfileImageUrl,
			},
		},
	})
	// Get twitter current information
}
func (h authHandler) DisconnectTwitterAccount(c *gin.Context) {
	authenticatedUser, err := h.app.Repositories.User.GetUserById(c.GetInt64(middlewares.KeyUserId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	if authenticatedUser.TwitterId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No twitter handle connected to this account",
		})
		return
	}

	user, err := h.app.Repositories.User.DisconnectTwitter(authenticatedUser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred while disconnecting user acccount",
			"err":     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Account disconnected successfully",
		"user":    user,
	})
}
