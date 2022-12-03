package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mypipeapp/mypipeapi/cmd/api/helpers"
	"github.com/mypipeapp/mypipeapi/cmd/api/internal"
	"github.com/mypipeapp/mypipeapi/cmd/api/middlewares"
	"github.com/mypipeapp/mypipeapi/db/actions/postgres"
	"github.com/mypipeapp/mypipeapi/db/models"
	"net/http"
)

type SearchHandler interface {
	Search(c *gin.Context)
}

type searchHandler struct {
	app internal.Application
}

func NewSearchHandler(app internal.Application) SearchHandler {
	return searchHandler{app: app}
}

func (h searchHandler) Search(c *gin.Context) {
	req := struct {
		Name string `form:"name"`
		Type string `form:"type"`
	}{}

	if err := c.Bind(&req); err != nil {
		errMessage := helpers.ParseErrorMessage(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": errMessage,
			"err":     err.Error(),
		})
		return
	}

	switch req.Type {
	case models.SearchTypePipes:
		pipes, err := h.app.Repositories.Search.SearchThroughPipes(req.Name, c.GetInt64(middlewares.KeyUserId))
		if err != nil {
			if err == postgres.ErrNoRecord {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"message": fmt.Sprintf("no results found for %v", req.Name),
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "an error occurred",
				"err":     err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"type": "pipe",
			"data": map[string]interface{}{
				"pipes": map[string]interface{}{
					"result": pipes,
					"total":  len(pipes),
				},
			},
		})

	case models.SearchTypeTags:
		bookmarks, err := h.app.Repositories.Search.SearchThroughTags(req.Name, c.GetInt64(middlewares.KeyUserId))
		if err != nil {
			if err == postgres.ErrNoRecord {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"message": fmt.Sprintf("no results found for %v", req.Name),
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "an error occurred",
				"err":     err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"type": "bookmark",
			"data": map[string]interface{}{
				"bookmarks": map[string]interface{}{
					"result": bookmarks,
					"total":  len(bookmarks),
				},
			},
		})

	case models.SearchTypePlatform:
		bookmarks, err := h.app.Repositories.Search.SearchThroughPlatform(req.Name, c.GetInt64(middlewares.KeyUserId))
		if err != nil {
			if err == postgres.ErrNoRecord {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"message": fmt.Sprintf("no results found for %v", req.Name),
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "an error occurred",
				"err":     err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"type": "bookmark",
			"data": map[string]interface{}{
				"bookmarks": map[string]interface{}{
					"result": bookmarks,
					"total":  len(bookmarks),
				},
			},
		})
	case models.SearchTypeAll:
		bookmarks, err := h.app.Repositories.Search.SearchThroughTags(req.Name, c.GetInt64(middlewares.KeyUserId))
		if err != nil {
			if err == postgres.ErrNoRecord {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"message": fmt.Sprintf("no results found for %v", req.Name),
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "an error occurred",
				"err":     err.Error(),
			})
			return
		}
		pipes, err := h.app.Repositories.Search.SearchThroughPipes(req.Name, c.GetInt64(middlewares.KeyUserId))
		if err != nil {
			if err == postgres.ErrNoRecord {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"message": fmt.Sprintf("no results found for %v", req.Name),
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "an error occurred",
				"err":     err.Error(),
			})
			return
		}

		platform, err := h.app.Repositories.Search.SearchThroughPlatform(req.Name, c.GetInt64(middlewares.KeyUserId))
		if err != nil {
			if err == postgres.ErrNoRecord {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"message": fmt.Sprintf("no results found for %v", req.Name),
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "an error occurred",
				"err":     err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"type": "all",
			"data": map[string]interface{}{
				"pipes": map[string]interface{}{
					"result": pipes,
					"total":  len(pipes),
				},
				"bookmarks": map[string]interface{}{
					"result": bookmarks,
					"total":  len(bookmarks),
				},
				"platform": map[string]interface{}{
					"result": platform,
					"total":  len(platform),
				},
				"total": len(bookmarks) + len(pipes) + len(platform),
			},
		})

	default:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid search type. valid search types are: *all*, *tags* and *pipes*",
		})
		return
	}
}
