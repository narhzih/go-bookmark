package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/mypipeapp/mypipeapi/cmd/api/internal"
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

}
