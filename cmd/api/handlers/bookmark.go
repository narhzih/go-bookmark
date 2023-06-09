package handlers

import (
	"github.com/mypipeapp/mypipeapi/cmd/api/internal"
	"github.com/mypipeapp/mypipeapi/cmd/api/middlewares"
	"github.com/mypipeapp/mypipeapi/db/actions/postgres"
	"github.com/mypipeapp/mypipeapi/db/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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
		Url   string  `json:"url" binding:"required"`
		Tags  string  `json:"tags"`
		Pipes []int64 `json:"pipes"`
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

	for _, pid := range bmRequest.Pipes {
		var err error
		if _, err = h.app.Services.UserOwnsPipe(pid, c.GetInt64(middlewares.KeyUserId)); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}
		bookmark = models.Bookmark{
			UserID:   c.GetInt64(middlewares.KeyUserId),
			PipeID:   pid,
			Platform: detectedPlatform,
			Url:      bmRequest.Url,
		}
		bookmark, err = h.app.Repositories.Bookmark.CreateBookmark(bookmark)
		if err != nil {
			if err == postgres.ErrRecordExists {
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

		// add tags to bookmark
		if strings.TrimSpace(bmRequest.Tags) != "" {
			parsedTags := strings.Split(bmRequest.Tags, ",")
			var tagsSlice []models.Tag
			for _, tagString := range parsedTags {
				tagData := models.Tag{Name: strings.TrimSpace(tagString)}
				tagsSlice = append(tagsSlice, tagData)
			}

			err = h.app.Repositories.Tag.AddTagsToBookmark(bookmark.ID, tagsSlice)
			if err != nil {
				h.app.Logger.Err(err).Msg("error occurred while trying to save tags")
			}
		}

		// don't really care if there's any error for now
		bookmark, _ = h.app.Repositories.Bookmark.ParseTags(bookmark)
	}

	// parse the tags as part of the bookmarks and send it back

	c.JSON(http.StatusCreated, gin.H{
		"message": "Url bookmarked successfully",
		"data": map[string]interface{}{
			"bookmark": map[string]interface{}{
				"id":        bookmark.ID,
				"url":       bookmark.Url,
				"platform":  bookmark.Platform,
				"tags":      bookmark.Tags,
				"createdAt": bookmark.CreatedAt,
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
	bookmark, err = h.app.Repositories.Bookmark.GetBookmark(bmId, c.GetInt64(middlewares.KeyUserId))
	if err != nil {
		if err == postgres.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Bookmark not found",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to retrieve bookmark",
			"err":     err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "bookmarked fetched successfully",
		"data": map[string]interface{}{
			"bookmark": map[string]interface{}{
				"id":        bookmark.ID,
				"url":       bookmark.Url,
				"platform":  bookmark.Platform,
				"createdAt": bookmark.CreatedAt,
				"tags":      bookmark.Tags,
			},
		},
	})
}
func (h bookmarkHandler) GetBookmarks(c *gin.Context) {
	userId := c.GetInt64(middlewares.KeyUserId)
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
	userId := c.GetInt64(middlewares.KeyUserId)
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
