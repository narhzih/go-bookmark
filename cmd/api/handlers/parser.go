package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/mypipeapp/mypipeapi/cmd/api/helpers"
	"github.com/mypipeapp/mypipeapi/cmd/api/internal"
	"github.com/mypipeapp/mypipeapi/db/models"
	"net/http"
)

type ParserHandler interface {
	TwitterLinkParser(c *gin.Context)
	GetCompleteThreadOfATweet(c *gin.Context)
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
		errMessage := helpers.ParseErrorMessage(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": errMessage,
			"err":     err.Error(),
		})
		return
	}

	parsedLink, err := helpers.ParseTwitterLink(reqBody.Link)
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
		errMessage := helpers.ParseErrorMessage(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": errMessage,
			"err":     err.Error(),
		})
		return
	}
	parsedLink, err := helpers.ParseYoutubeLink(reqBody.Link)
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
		errMessage := helpers.ParseErrorMessage(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": errMessage,
			"err":     err.Error(),
		})
		return
	}
	parsedLink, err := helpers.ParseLink(reqBody.Link)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "An error occurred while parsing youtube link",
			"err":     err.Error(),
		})
		return
	}
	c.Data(http.StatusOK, "application/json", []byte(parsedLink))
}

func (h parserHandler) GetCompleteThreadOfATweet(c *gin.Context) {
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
	chatId := helpers.GetTwitterChatId(reqBody.Link)
	expandedResponse, err := h.app.Services.ExpandTweet(chatId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "an error occurred while connecting to twitter api",
			"err":     err.Error(),
		})
		return
	}

	expandedData := expandedResponse.Data[0]
	authorInfo, err := h.app.Services.GetFullUserInformation(expandedData.AuthorID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "an error occurred while connecting to twitter api",
			"err":     err.Error(),
		})
		return
	}
	completeTweet, err := h.app.Services.GetThreadByConversationID(expandedData.ConversationID, authorInfo.Data.Username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "an error occurred while connecting to twitter api",
			"err":     err.Error(),
		})
		return
	}
	var thread models.TwitterExpandedResponse
	json.Unmarshal([]byte(completeTweet), &thread)
	thread.Data = append(thread.Data, expandedData)
	c.JSON(http.StatusOK, gin.H{
		"thread": thread,
		"author": authorInfo,
	})

	//c.Data(http.StatusOK, "application/json", []byte(completeTweet))
}
