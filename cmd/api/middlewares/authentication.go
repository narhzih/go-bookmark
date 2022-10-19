package middlewares

import (
	"fmt"
	"gitlab.com/trencetech/mypipe-api/cmd/api/internal"
	"gitlab.com/trencetech/mypipe-api/db/actions/postgres"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const (
	KeyUserId   = "user_id"
	KeyUsername = "username"
)

var (
	InvalidToken = fmt.Errorf("no token present in request")
)

func AuthRequired(app internal.Application, jwtSecret string) gin.HandlerFunc {

	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader != "" {
			authToken := strings.Split(authHeader, " ")[1]
			token, err := jwt.Parse(authToken, func(t *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil {
				app.Logger.Err(err).Msg(err.Error())
				if err == InvalidToken {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"message": "You are not logged in. Please log in!",
					})
					return
				}
				app.Logger.Info().Msg("Error is from here")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": fmt.Sprintf("could not parse authorization token: %s", err),
				})
				return
			}

			claims := token.Claims.(jwt.MapClaims)
			username := claims["username"].(string)
			userId := int64(claims["sub"].(float64))
			app.Logger.Info().Msg(fmt.Sprintf("parsed username as: %v and userId as %v", username, userId))

			if username == "" || userId == 0 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "invalid user data in provided token",
				})
				return
			}

			// see if user still exists
			loggedInUser, err := app.Repositories.User.GetUserById(int(userId))
			if err != nil {
				if err == postgres.ErrNoRecord {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"message": "Session expired. Please login again",
						"err":     err.Error(),
					})
					return
				}
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Authentication error",
					"err":     err.Error(),
				})
				return
			}

			if loggedInUser.Username == username && loggedInUser.ID == userId {
				c.Set(KeyUsername, username)
				c.Set(KeyUserId, userId)
				c.Next()
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "Unauthorized",
				})
				return
			}

		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "You're not logged in!. Please login to perform this operation.",
			})
			return
		}

	}
}
