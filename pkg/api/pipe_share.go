package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/trencetech/mypipe-api/db"
	"net/http"
	"strconv"
)

func (h *Handler) SharePipe(c *gin.Context) {
	shareType := c.Query("type")
	if len(shareType) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Please specify how you want to share your pipe",
		})
		return
	}

	loggedInUser := c.GetInt64(KeyUserId)
	pipeID, _ := strconv.Atoi(c.Param("id"))
	_, err := h.service.DB.GetPipe(int64(pipeID), loggedInUser)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Our system encountered an error while trying to finalize share. Please try again soon",
		})
		return
	}

	if shareType == "public" {
		sharedPipe, err := h.service.SharePublicPipe(int64(pipeID), loggedInUser)
		if err != nil {
			h.logger.Err(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Our system encountered an error while trying to finalize share. Please try again soon",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Public pipe share created successfully. Get the code and share it to whomever you want to share your pipe with",
			"data": map[string]interface{}{
				"share_code": sharedPipe.Code,
			},
		})
	} else if shareType == "private" {
		shareReq := struct {
			Username string `json:"username" binding:"required"`
		}{}
		if err := c.ShouldBindJSON(&shareReq); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Please provide a user to share pipe with",
			})
		}
		shareTo := shareReq.Username
		_, err := h.service.SharePrivatePipe(int64(pipeID), loggedInUser, shareTo)
		if err != nil {
			if err == db.ErrPipeShareToNotFound || err == db.ErrCannotSharePipeToSelf {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
				})
				return
			}
			h.logger.Err(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Our system encountered an error while trying to finalize share. Please try again soon",
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message": fmt.Sprintf("Pipe has been shared to %v successfully", shareTo),
		})
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid pipe share type",
		})
	}
}
func (h *Handler) ReceivePipe(c *gin.Context)               {}
func (h *Handler) RemoveShareAccessFromPipe(c *gin.Context) {}
func (h *Handler) ChangePipeShareAccessType(c *gin.Context) {}
