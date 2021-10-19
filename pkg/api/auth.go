package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SignIn(c *gin.Context) {
	c.JSON(http.StatusAccepted, gin.H{
		"message": "All good",
	})
}

func (h *Handler) SingUp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Implement route to sign user up",
	})
}
