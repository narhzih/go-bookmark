package api

import "github.com/gin-gonic/gin"

func (h *Handler) SharePipe(c *gin.Context)                 {}
func (h *Handler) ReceivePipe(c *gin.Context)               {}
func (h *Handler) RemoveShareAccessFromPipe(c *gin.Context) {}
func (h *Handler) ChangePipeShareAccessType(c *gin.Context) {}
