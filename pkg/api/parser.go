package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/trencetech/mypipe-api/pkg/helpers"
	"gitlab.com/trencetech/mypipe-api/pkg/service/mp_parser"
	"net/http"
)

func (h *Handler) TwitterLinkParser(c *gin.Context) {
	reqBody := struct {
		Link string `json:"link"`
	}{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		errMessage := helpers.ParseErrorMessage(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": errMessage,
			"err":     err.Error(),
		})
		return
	}

	parsedLink, err := mp_parser.ParseTwitterLink(reqBody.Link)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "An error occurred while parsing twitter link",
			"err":     err.Error(),
		})
		return
	}
	c.Data(http.StatusOK, "application/json", []byte(parsedLink))
}