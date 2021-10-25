package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/request"
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
		token, err := request.ParseFromRequest(c.Request, request.AuthorizationHeaderExtractor, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil {
			if err == InvalidToken {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "You are not logged in. Please log in!",
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": fmt.Sprintf("could not parse authorization token: %s", err),
			})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		username := claims["username"].(string)
		userId := int(claims["sub"].(float64))

		logger.Info().Msgf("%+v", claims)
		if username == "" || userId == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid user data in provided token",
			})
			return
		}

		c.Set(KeyUsername, username)
		c.Set(KeyUserId, userId)
		c.Next()
	}
}
