package handlers

import (
	"gitlab.com/trencetech/mypipe-api/cmd/api/internal"
	"gitlab.com/trencetech/mypipe-api/db/actions/postgres"
	"gitlab.com/trencetech/mypipe-api/db/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/trencetech/mypipe-api/db"
)

type BookmarkHandler interface {
	GetBookmarks(c *gin.Context)
	CreateBookmark(c *gin.Context)
	GetBookmark(c *gin.Context)
	DeleteBookmark(c *gin.Context)
}

type bookmarkHandler struct {
	app internal.Application
}

func NewBookmarkHandler(app internal.Application) BookmarkHandler {
	return bookmarkHandler{app: app}
}

func (h bookmarkHandler) CreateBookmark(c *gin.Context) {
	bmRequest := struct {
		Url string `json:"url" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&bmRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}
	var detectedPlatform string
	detectedPlatform, _ = h.app.Services.GetPlatformFromLink(bmRequest.Url)
	var bookmark models.Bookmark
	pipeId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.app.Logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid pipe ID",
		})
		return
	}
	pipeExits, err := h.app.Services.PipeExists(pipeId, c.GetInt64(KeyUserId))
	if !pipeExits {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid pipe to create bookmark",
		})
		return
	}
	if _, err := h.app.Services.UserOwnsPipe(pipeId, c.GetInt64(KeyUserId)); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
	}
	bookmark = models.Bookmark{
		UserID:   c.GetInt64(KeyUserId),
		PipeID:   pipeId,
		Platform: detectedPlatform,
		Url:      bmRequest.Url,
	}
	bookmark, err = h.app.Repositories.Bookmark.CreateBookmark(bookmark)
	if err != nil {
		if err == db.ErrRecordExists {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "You have already bookmarked this url",
			})
			return
		}
		h.app.Logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred when trying to create your bookmark",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Url bookmarked successfully",
		"data": map[string]interface{}{
			"bookmark": map[string]interface{}{
				"id":       bookmark.ID,
				"url":      bookmark.Url,
				"platform": bookmark.Platform,
			},
		},
	})
}
func (h bookmarkHandler) GetBookmark(c *gin.Context) {
	var bookmark models.Bookmark
	bmId, err := strconv.ParseInt(c.Param("bmId"), 10, 64)
	if err != nil {
		h.app.Logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Bookmark ID",
		})
		return
	}
	bookmark, err = h.app.Repositories.Bookmark.GetBookmark(bmId, c.GetInt64(KeyUserId))
	if err != nil {
		if err == postgres.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Bookmark not found",
			})
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to retreive bookmark",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "bookmarked fetched successfully",
		"data": map[string]interface{}{
			"bookmark": map[string]interface{}{
				"id":       bookmark.ID,
				"url":      bookmark.Url,
				"platform": bookmark.Platform,
			},
		},
	})
}
func (h bookmarkHandler) GetBookmarks(c *gin.Context) {
	userId := c.GetInt64(KeyUserId)
	pipeId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.app.Logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Bookmark ID",
		})
		return
	}
	bookmarks, err := h.app.Repositories.Bookmark.GetBookmarks(userId, pipeId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Could not retrieve bookmarks! Please try again soon",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pipes fetched successfully",
		"data": map[string]interface{}{
			"bookmarks": bookmarks,
		},
	})
}
func (h bookmarkHandler) DeleteBookmark(c *gin.Context) {
	userId := c.GetInt64(KeyUserId)
	bmId, err := strconv.ParseInt(c.Param("bmId"), 10, 64)
	if err != nil {
		h.app.Logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Bookmark ID",
		})
		return
	}

	_, err = h.app.Repositories.Bookmark.DeleteBookmark(bmId, userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred while trying to delete bookmark",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bookmark deleted successfully",
	})
}
