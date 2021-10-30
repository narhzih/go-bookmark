package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog"
)

const (
	KeyUserId   = "user_id"
	KeyUsername = "username"
)

var (
	InvalidToken = fmt.Errorf("no token present in request")
)

func AuthRequired(jwtSecret string, logger zerolog.Logger) gin.HandlerFunc {

	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader != "" {
			authToken := strings.Split(authHeader, " ")[1]
			token, err := jwt.Parse(authToken, func(t *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil {
				logger.Err(err).Msg(err.Error())
				if err == InvalidToken {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"message": "You are not logged in. Please log in!",
					})
					return
				}
				logger.Info().Msg("Error is from here")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": fmt.Sprintf("could not parse authorization token: %s", err),
				})
				return
			}

			claims := token.Claims.(jwt.MapClaims)
			username := claims["username"].(string)
			userId := int64(claims["sub"].(float64))

			if username == "" || userId == 0 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "invalid user data in provided token",
				})
				return
			}

			c.Set(KeyUsername, username)
			c.Set(KeyUserId, userId)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "You're not logged in!. Please login to perform this operation.",
			})
			return
		}

	}
}
