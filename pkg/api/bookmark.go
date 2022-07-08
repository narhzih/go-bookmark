package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/trencetech/mypipe-api/db"
	"gitlab.com/trencetech/mypipe-api/db/model"
)

func (h *Handler) CreateBookmark(c *gin.Context) {
	bmRequest := struct {
		Platform string `json:"platform" binding:"required"`
		Url      string `json:"url" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&bmRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	var bookmark model.Bookmark
	pipeId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid pipe ID",
		})
		return
	}
	pipeExits, err := h.service.PipeExists(pipeId, c.GetInt64(KeyUserId))
	if !pipeExits {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid pipe to create bookmark",
		})
		return
	}
	if _, err := h.service.UserOwnsPipe(pipeId, c.GetInt64(KeyUserId)); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
	}
	bookmark = model.Bookmark{
		UserID:   c.GetInt64(KeyUserId),
		PipeID:   pipeId,
		Platform: bmRequest.Platform,
		Url:      bmRequest.Url,
	}
	bookmark, err = h.service.DB.CreateBookmark(bookmark)
	if err != nil {
		if err == db.ErrRecordExists {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "You have already bookmarked this url",
			})
			return
		}
		h.logger.Err(err).Msg(err.Error())
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
func (h *Handler) GetBookmark(c *gin.Context) {
	var bookmark model.Bookmark
	bmId, err := strconv.ParseInt(c.Param("bmId"), 10, 64)
	if err != nil {
		h.logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Bookmark ID",
		})
		return
	}
	bookmark, err = h.service.DB.GetBookmark(bmId, c.GetInt64(KeyUserId))
	if err != nil {
		if err == db.ErrNoRecord {
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
func (h *Handler) GetBookmarks(c *gin.Context) {
	userId := c.GetInt64(KeyUserId)
	pipeId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Bookmark ID",
		})
		return
	}
	bookmarks, err := h.service.DB.GetBookmarks(userId, pipeId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Could not retreive bookmarks! Please try again soon",
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
func (h *Handler) DeleteBookmark(c *gin.Context) {
	userId := c.GetInt64(KeyUserId)
	bmId, err := strconv.ParseInt(c.Param("bmId"), 10, 64)
	if err != nil {
		h.logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Bookmark ID",
		})
		return
	}

	_, err = h.service.DB.DeleteBookmark(bmId, userId)
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
