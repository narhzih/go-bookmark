package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/trencetech/mypipe-api/db"
	"gitlab.com/trencetech/mypipe-api/db/model"
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

func (h *Handler) PreviewPipe(c *gin.Context) {
	code := c.Query("code")
	// See if the pipe is still sharable
	pipeToAdd, err := h.service.DB.GetSharedPipeByCode(code)
	if err != nil {
		if err == db.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "This pipe cannot be added to your collection because the author has not allowed it!",
			})
			return
		}
		h.logger.Err(err).Msg("Something went wrong")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Our system encountered an error while trying to add pipe your collection. Try again soon!",
		})
		return
	}
	pipe, err := h.service.DB.GetPipeAndResource(pipeToAdd.PipeID, pipeToAdd.SharerID)
	if err != nil {
		if err == db.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Pipe not found. Either the author of this pipe has deleted it or has removed the pipe from public access",
			})
			return
		}
		h.logger.Err(err).Msg("Something went wrong")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Our system encountered an error while trying to add pipe your collection. Try again soon!",
		})
		return
	}
	sharer, err := h.service.DB.GetUserById(int(pipeToAdd.SharerID))
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"err":     err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Preview successful",
		"data": map[string]interface{}{
			"user": map[string]interface{}{
				"username": sharer.Username,
				"name":     sharer.ProfileName,
			},
			"pipe": pipe,
		},
	})

}

func (h *Handler) AddPipe(c *gin.Context) {
	code := c.Query("code")
	// See if the pipe is still sharable
	pipeToAdd, err := h.service.DB.GetSharedPipeByCode(code)
	if err != nil {
		if err == db.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "This pipe cannot be added to your collection because the author has not allowed it!",
			})
			return
		}
		h.logger.Err(err).Msg("Something went wrong")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Our system encountered an error while trying to add pipe your collection. Try again soon!",
		})
		return
	}
	h.logger.Info().Msg("Pipe was not found and about to be added to users collection")
	// See if this user has already added this pipe to their collection
	_, err = h.service.DB.GetReceivedPipeRecord(pipeToAdd.ID, c.GetInt64(KeyUserId))
	if err == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "You have already added this pipe to your collection.",
		})
		return
	} else {
		if err != db.ErrNoRecord {
			// This
			h.logger.Err(err).Msg("An error occurred")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "An error occurred while trying to add pipe to collection.",
			})
			return
		}
	}
	_, err = h.service.DB.GetPipe(pipeToAdd.PipeID, pipeToAdd.SharerID)
	if err != nil {
		if err == db.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Pipe not found. Either the author of this pipe has deleted it or has removed the pipe from public access",
			})
			return
		}
		h.logger.Err(err).Msg("Something went wrong")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Our system encountered an error while trying to add pipe your collection. Try again soon!",
		})
		return
	}

	if pipeToAdd.Type != "public" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "You cannot add this pipe to your collection because this pipe is not a public pipe",
		})
		return
	}

	if pipeToAdd.SharerID == c.GetInt64(KeyUserId) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "You cannot add pipe to collection because it's yours!",
		})
	}

	// Now we can add the pipe to the user's collection
	_, err = h.service.DB.CreatePipeReceiver(model.SharedPipeReceiver{
		SharedPipeId: pipeToAdd.ID,
		ReceiverID:   c.GetInt64(KeyUserId),
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"err":     err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pipe has been added to your collection successfully",
	})
}
func (h *Handler) RemoveShareAccessFromPipe(c *gin.Context) {}
func (h *Handler) ChangePipeShareAccessType(c *gin.Context) {}
