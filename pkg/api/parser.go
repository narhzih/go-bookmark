package api

import (
	"github.com/gin-gonic/gin"
	mp_parser2 "gitlab.com/trencetech/mypipe-api/cmd/api/services/mp_parser"
	"gitlab.com/trencetech/mypipe-api/pkg/helpers"
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

	parsedLink, err := mp_parser2.ParseTwitterLink(reqBody.Link)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "An error occurred while parsing twitter link",
			"err":     err.Error(),
		})
		return
	}
	c.Data(http.StatusOK, "application/json", []byte(parsedLink))
}

func (h *Handler) YoutubeLinkParser(c *gin.Context) {
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
	parsedLink, err := mp_parser2.ParseYoutubeLink(reqBody.Link)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "An error occurred while parsing youtube link",
			"err":     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": parsedLink,
	})
}

func (h *Handler) ParseLink(c *gin.Context) {
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
	parsedLink, err := mp_parser2.ParseLink(reqBody.Link)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "An error occurred while parsing youtube link",
			"err":     err.Error(),
		})
		return
	}
	c.Data(http.StatusOK, "application/json", []byte(parsedLink))
}
