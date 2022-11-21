package middlewares

import "github.com/gin-gonic/gin"

type AuthenticatedUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

func GetLoggedInUser(c *gin.Context) AuthenticatedUser {
	return AuthenticatedUser{
		ID:       c.GetInt64(KeyUserId),
		Username: c.GetString(KeyUsername),
	}
}
