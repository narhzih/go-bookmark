package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/trencetech/mypipe-api/cmd/api/internal"
	"gitlab.com/trencetech/mypipe-api/db"
	"gitlab.com/trencetech/mypipe-api/db/models"
	"net/http"
	"strconv"
)

type PipeShareHandler interface {
	SharePipe(c *gin.Context)
	PreviewPipe(c *gin.Context)
	AddPipe(c *gin.Context)
}

type pipeShareHandler struct {
	app internal.Application
}

func NewPipeShareHandler(app internal.Application) PipeShareHandler {
	return pipeShareHandler{app: app}
}

func (h pipeShareHandler) SharePipe(c *gin.Context) {
	shareType := c.Query("type")
	if len(shareType) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Please specify how you want to share your pipe",
		})
		return
	}

	loggedInUser := c.GetInt64(KeyUserId)
	pipeID, _ := strconv.Atoi(c.Param("id"))
	_, err := h.app.Repositories.Pipe.GetPipe(int64(pipeID), loggedInUser)
	if err != nil {
		h.app.Logger.Err(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Our system encountered an error while trying to finalize share. Please try again soon",
		})
		return
	}

	if shareType == "public" {
		sharedPipe, err := h.app.Services.SharePublicPipe(int64(pipeID), loggedInUser)
		if err != nil {
			h.app.Logger.Err(err)
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
		userToShareTo, err := h.app.Repositories.User.GetUserByUsername(shareTo)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "The specified user not found",
			})
			return
		}
		_, err = h.app.Services.SharePrivatePipe(int64(pipeID), loggedInUser, shareTo)
		if err != nil {
			if err == db.ErrPipeShareToNotFound || err == db.ErrCannotSharePipeToSelf {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
				})
				return
			}
			h.app.Logger.Err(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Our system encountered an error while trying to finalize share. Please try again soon",
			})
			return
		}

		err = h.app.Services.CreatePrivatePipeShareNotification(int64(pipeID), loggedInUser, userToShareTo.ID)
		if err != nil {
			h.app.Logger.Err(err).Msg("An error occurred while creating share notification")
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

func (h pipeShareHandler) PreviewPipe(c *gin.Context) {
	code := c.Query("code")
	// See if the pipe is still sharable
	pipeToAdd, err := h.app.Repositories.PipeShare.GetSharedPipeByCode(code)
	if err != nil {
		if err == db.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "This pipe cannot be added to your collection because the author has not allowed it!",
			})
			return
		}
		h.app.Logger.Err(err).Msg("Something went wrong")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Our system encountered an error while trying to add pipe your collection. Try again soon!",
		})
		return
	}
	pipe, err := h.app.Repositories.Pipe.GetPipeAndResource(pipeToAdd.PipeID, pipeToAdd.SharerID)
	if err != nil {
		if err == db.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Pipe not found. Either the author of this pipe has deleted it or has removed the pipe from public access",
			})
			return
		}
		h.app.Logger.Err(err).Msg("Something went wrong")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Our system encountered an error while trying to add pipe your collection. Try again soon!",
		})
		return
	}
	sharer, err := h.app.Repositories.User.GetUserById(int(pipeToAdd.SharerID))
	if err != nil {
		h.app.Logger.Err(err)
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

func (h pipeShareHandler) AddPipe(c *gin.Context) {
	pipeId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.app.Logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid pipe ID",
		})
		return
	}

	// See if the pipe is an actual pipe
	pipeToAdd, err := h.app.Repositories.PipeShare.GetSharedPipe(pipeId)
	if err != nil {
		switch {
		case errors.Is(err, db.ErrNoRecord):
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "This pipe cannot be added to your collection because the author has not allowed it!",
			})
			return
		default:
			h.app.Logger.Err(err).Msg("Something went wrong")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Our system encountered an error while trying to add pipe your collection. Try again soon!",
			})
			return

		}

	}

	if pipeToAdd.SharerID == c.GetInt64(KeyUserId) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "You cannot add pipe to collection because it's yours!",
		})
	}

	switch pipeToAdd.Type {
	case models.PipeShareTypePublic:
		// Perform actions for public share type
		_, err := h.app.Repositories.PipeShare.GetSharedPipeByCode(c.Query("code"))
		if err != nil {
			switch {
			case errors.Is(err, db.ErrNoRecord):
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": "This pipe cannot be added to your collection because the author has not allowed it!",
				})
				return
			default:
				h.app.Logger.Err(err).Msg("Something went wrong")
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Our system encountered an error while trying to add pipe your collection. Try again soon!",
				})
				return

			}
		}

		// See if this user has already added this pipe to their collection
		_, err = h.app.Repositories.PipeShare.GetReceivedPipeRecord(pipeToAdd.ID, c.GetInt64(KeyUserId))
		switch {
		case errors.Is(err, nil):
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "You have already added this pipe to your collection.",
			})
			return
		case err != db.ErrNoRecord:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "An error occurred while trying to add pipe to collection.",
			})
			return
		}

		// get the actual pipe
		_, err = h.app.Repositories.Pipe.GetPipe(pipeToAdd.PipeID, pipeToAdd.SharerID)
		if err != nil {
			if err == db.ErrNoRecord {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": "Pipe not found. Either the author of this pipe has deleted it or has removed the pipe from public access",
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Our system encountered an error while trying to add pipe your collection. Try again soon!",
			})
			return
		}

		_, err = h.app.Repositories.PipeShare.CreatePipeReceiver(models.SharedPipeReceiver{
			SharedPipeId: pipeToAdd.ID,
			ReceiverID:   c.GetInt64(KeyUserId),
			IsAccepted:   true,
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong",
				"err":     err.Error(),
			})
			return
		}
	case models.PipeShareTypePrivate:
		pipeShareRecord, err := h.app.Repositories.PipeShare.GetReceivedPipeRecord(pipeToAdd.PipeID, c.GetInt64(KeyUserId))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "an error occurred while addding pipe to your collection",
				"err":     err.Error(),
			})
			return
		}

		if pipeShareRecord.IsAccepted {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"message": "this pipe is already in your collection",
			})
		}
		_, err = h.app.Repositories.PipeShare.AcceptPrivateShare(pipeShareRecord)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "an error occurred while addding pipe to your collection",
				"err":     err.Error(),
			})
			return
		}

	default:
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message": "Operation not allowed for pipe sharing",
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pipe has been added to your collection successfully",
	})
}

func (h pipeShareHandler) RemoveShareAccessFromPipe(c *gin.Context) {}
func (h pipeShareHandler) ChangePipeShareAccessType(c *gin.Context) {}
