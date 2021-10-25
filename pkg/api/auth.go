package api

import (
	"net/http"

	"gitlab.com/gowagr/mypipe-api/db"
	"gitlab.com/gowagr/mypipe-api/db/model"

	"github.com/gin-gonic/gin"
)

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

func (h *Handler) SingUpWithGoogle(c *gin.Context) {
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
				"message": "Email has already been taken.",
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
