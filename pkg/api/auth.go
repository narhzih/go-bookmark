package api

import (
	"net/http"

	"gitlab.com/trencetech/mypipe-api/db"
	"gitlab.com/trencetech/mypipe-api/db/model"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

func (h *Handler) EmailSignUp(c *gin.Context) {
	singUpReq := struct {
		Username    string `json:"username" binding:"required"`
		ProfileName string `json:"profile_name" binding:"required"`
		Email       string `json:"email" binding:"required"`
		Password    string `json:"password" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&singUpReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	hashedPassword, err := hashPassword(singUpReq.Password)
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

	_, err = h.service.DB.CreateUserByEmail(userStruct, hashedPassword)
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
	/**
	* TODO @narhzih
	* Implement verification email step after registration
	 */

	c.JSON(http.StatusOK, gin.H{
		"message": "Account created successfully. Please check your email for verification code",
	})

	// authToken, err := h.service.IssueAuthToken(user)
	// if err != nil {
	// 	h.logger.Err(err).Msg(err.Error())
	// 	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
	// 		"message": "Error occurred while trying to sign user in",
	// 	})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "Sign up successful!",
	// 	"data": map[string]interface{}{
	// 		"token":         authToken.AccessToken,
	// 		"refresh_token": authToken.RefreshToken,
	// 		"user": map[string]interface{}{
	// 			"id":       user.ID,
	// 			"username": user.Username,
	// 			"email":    user.Email,
	// 		},
	// 	},
	// })

}

func (h *Handler) EmailLogin(c *gin.Context) {
	loginReq := struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		h.logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
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
		if err == db.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "You cannot be logged in! Your account was either created through google sign up or apple sign up. Please use either of those to sign in to your account",
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred while trying to sign in user",
		})
		return
	}

	if verifyPassword(loginReq.Password, userAndAuth.HashedPassword) {
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
				"user": map[string]interface{}{
					"id":       user.ID,
					"username": user.Username,
					"email":    user.Email,
				},
			},
		})
	} else {
		h.logger.Err(err).Msg(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Password incorrect",
		})
		return
	}

	// Validate the password provided by the user

}

func (h *Handler) SignInWithGoogle(c *gin.Context) {
	var user model.User

	singInReq := struct {
		TokenString string `json:"token_string" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&singInReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	claims, err := h.service.ValidateGoogleJWT(singInReq.TokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid GoogleJWT",
		})
		return
	}

	user, err = h.service.DB.GetUserByEmail(claims.Email)
	if err != nil {
		if err == db.ErrNoRecord {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "user with the specified email not found",
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
			"message": "Error occurred while trying to sign user in",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sign in successful!",
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
	claims, err := h.service.ValidateGoogleJWT(signUpReq.TokenString)
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
			"token":        authToken.AccessToken,
			"refesh_token": authToken.RefreshToken,
			"user": map[string]interface{}{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
			},
		},
	})
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
