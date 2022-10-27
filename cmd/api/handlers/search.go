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
		Name string `query:"name"`
		Type string `query:"type"`
	}{}

	if err := c.ShouldBindJSON(&req); err != nil {
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
			"query": req.Type,
			"result": map[string]interface{}{
				"pipes": pipes,
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
			"query": req.Type,
			"result": map[string]interface{}{
				"bookmarks": bookmarks,
			},
		})

	case models.SearchTypeAll:
		panic("implement search for all (both tags and pipes)")
	default:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid search type. valid search types are: *all*, *tags* and *pipes*",
		})
		return
	}
}
