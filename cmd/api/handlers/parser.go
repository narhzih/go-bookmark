package handlers

import (
	"github.com/gin-gonic/gin"
	helpers2 "github.com/mypipeapp/mypipeapi/cmd/api/helpers"
	"github.com/mypipeapp/mypipeapi/cmd/api/internal"
	"net/http"
)

type ParserHandler interface {
	TwitterLinkParser(c *gin.Context)
	YoutubeLinkParser(c *gin.Context)
	ParseLink(c *gin.Context)
}

type parserHandler struct {
	app internal.Application
}

func NewParserHandler(app internal.Application) ParserHandler {
	return parserHandler{app: app}
}

func (h parserHandler) TwitterLinkParser(c *gin.Context) {
	reqBody := struct {
		Link string `json:"link"`
	}{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		errMessage := helpers2.ParseErrorMessage(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": errMessage,
			"err":     err.Error(),
		})
		return
	}

	parsedLink, err := helpers2.ParseTwitterLink(reqBody.Link)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "An error occurred while parsing twitter link",
			"err":     err.Error(),
		})
		return
	}
	c.Data(http.StatusOK, "application/json", []byte(parsedLink))
}

func (h parserHandler) YoutubeLinkParser(c *gin.Context) {
	reqBody := struct {
		Link string `json:"link"`
	}{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		errMessage := helpers2.ParseErrorMessage(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": errMessage,
			"err":     err.Error(),
		})
		return
	}
	parsedLink, err := helpers2.ParseYoutubeLink(reqBody.Link)
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

func (h parserHandler) ParseLink(c *gin.Context) {
	reqBody := struct {
		Link string `json:"link"`
	}{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		errMessage := helpers2.ParseErrorMessage(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": errMessage,
			"err":     err.Error(),
		})
		return
	}
	parsedLink, err := helpers2.ParseLink(reqBody.Link)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "An error occurred while parsing youtube link",
			"err":     err.Error(),
		})
		return
	}
	c.Data(http.StatusOK, "application/json", []byte(parsedLink))
}
