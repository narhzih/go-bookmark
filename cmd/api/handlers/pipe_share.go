package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mypipeapp/mypipeapi/cmd/api/internal"
	"github.com/mypipeapp/mypipeapi/cmd/api/middlewares"
	"github.com/mypipeapp/mypipeapi/db/actions/postgres"
	"github.com/mypipeapp/mypipeapi/db/models"
	"net/http"
	"strconv"
	"strings"
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
	reqType := c.Query("type")
	if strings.TrimSpace(reqType) == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Please specify a valid share type",
		})
		return
	}

	// Validate inputs integrity
	req := struct {
		Username string `form:"username" json:"username"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	h.app.Logger.Info().Msg(fmt.Sprintf("username -> %s", req.Username))

	sharerId := c.GetInt64(middlewares.KeyUserId)
	id, _ := strconv.Atoi(c.Param("id"))
	pipeId := int64(id)
	_, err := h.app.Repositories.Pipe.GetPipe(pipeId, sharerId)
	if err != nil {
		h.app.Logger.Err(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid pipe ID",
		})
		return
	}

	switch reqType {
	case models.PipeShareTypePublic:
		sharedPipe, err := h.app.Services.SharePipePublicly(pipeId, sharerId)
		if err != nil {
			h.app.Logger.Err(err).Msg("an error occurred while sharing pipe publicly")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Our system encountered an error while trying to finalize share. Please try again soon",
				"err":     err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Public pipe share created successfully. Get the code and share it to whomever you want to share your pipe with",
			"data": map[string]interface{}{
				"share_code": sharedPipe.Code,
			},
		})
		return

	case models.PipeShareTypePrivate:
		// validate that a valid username was sent alongside the request
		if strings.TrimSpace(req.Username) == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "For a private pipe share, you have to specify a username you want to share to",
			})
			return
		}
		receiver, err := h.app.Repositories.User.GetUserByUsername(req.Username)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "The specified user not found",
			})
			return
		}
		_, err = h.app.Services.SharePipePrivately(pipeId, sharerId, receiver.Username)
		if err != nil {
			h.app.Logger.Err(err).Msg("an error occurred while sharing pipe publicly")
			if err == postgres.ErrPipeShareToNotFound || err == postgres.ErrCannotSharePipeToSelf {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Our system encountered an error while trying to finalize share. Please try again soon",
				"err":     err.Error(),
			})
			return
		}

		err = h.app.Services.CreatePrivatePipeShareNotification(pipeId, sharerId, receiver.ID)
		if err != nil {
			h.app.Logger.Err(err).Msg("An error occurred while creating share notification")
		}
		c.JSON(http.StatusCreated, gin.H{
			"message": fmt.Sprintf("Pipe has been shared with %v successfully", receiver.Username),
		})

	default:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Share types can either be public or private",
		})
		return
	}

}

func (h pipeShareHandler) PreviewPipe(c *gin.Context) {
	code := c.Query("code")
	// See if the pipe is still sharable
	authenticatedUser, _ := h.app.Repositories.User.GetUserById(c.GetInt64(middlewares.KeyUserId))

	pipeToAdd, err := h.app.Repositories.PipeShare.GetSharedPipeByCode(code)
	if err != nil {
		if err == postgres.ErrNoRecord {
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

	// If it's a private pipe, check if this user is allowed to see it
	if pipeToAdd.Type == models.PipeShareTypePrivate {
		_, err := h.app.Repositories.PipeShare.GetReceivedPipeRecord(pipeToAdd.ID, authenticatedUser.ID)
		if err != nil {
			if err == postgres.ErrNoRecord {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "You cannot view this pipe because it's a private pipe",
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Our server encountered an error. Please try again in few minutes",
				"err":     err.Error(),
			})
			return
		}
	}

	pipe, err := h.app.Repositories.Pipe.GetPipeAndResource(pipeToAdd.PipeID, pipeToAdd.SharerID)
	if err != nil {
		if err == postgres.ErrNoRecord {
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
	sharer, err := h.app.Repositories.User.GetUserById(pipeToAdd.SharerID)
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
			"sharer": map[string]interface{}{
				"username":    sharer.Username,
				"name":        sharer.ProfileName,
				"cover_photo": sharer.CovertPhoto,
			},
			"fullPipeData": pipe,
		},
	})

}

func (h pipeShareHandler) AddPipe(c *gin.Context) {
	code := c.Query("code")

	// See if the pipe is an actual pipe
	pipeToAdd, err := h.app.Repositories.PipeShare.GetSharedPipeByCode(code)
	if err != nil {
		switch {
		case errors.Is(err, postgres.ErrNoRecord):
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

	if pipeToAdd.SharerID == c.GetInt64(middlewares.KeyUserId) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "You cannot add pipe to collection because it's yours!",
		})
		return
	}

	h.app.Logger.Info().Msg("Pipe share type to add is ->" + pipeToAdd.Type)
	switch pipeToAdd.Type {
	case models.PipeShareTypePublic:
		h.app.Logger.Info().Msg("adding pipe publicly")
		// Perform actions for public share type
		_, err := h.app.Repositories.PipeShare.GetSharedPipeByCode(c.Query("code"))
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNoRecord):
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
		_, err = h.app.Repositories.PipeShare.GetReceivedPipeRecord(pipeToAdd.ID, c.GetInt64(middlewares.KeyUserId))
		switch {
		case errors.Is(err, nil):
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "You have already added this pipe to your collection.",
			})
			return
		case err != postgres.ErrNoRecord:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "An error occurred while trying to add pipe to collection.",
			})
			return
		}

		// get the actual pipe
		_, err = h.app.Repositories.Pipe.GetPipe(pipeToAdd.PipeID, pipeToAdd.SharerID)
		if err != nil {
			if err == postgres.ErrNoRecord {
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
			ReceiverID:   c.GetInt64(middlewares.KeyUserId),
			IsAccepted:   true,
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
		return
	case models.PipeShareTypePrivate:
		h.app.Logger.Info().Msg("adding pipe privately")
		pipeShareRecord, err := h.app.Repositories.PipeShare.GetReceivedPipeRecord(pipeToAdd.PipeID, c.GetInt64(middlewares.KeyUserId))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "an error occurred while adding pipe to your collection",
				"err":     err.Error(),
			})
			return
		}

		if pipeShareRecord.IsAccepted {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"message": "this pipe is already in your collection",
			})
			return
		}
		_, err = h.app.Repositories.PipeShare.AcceptPrivateShare(pipeShareRecord)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "an error occurred while adding pipe to your collection",
				"err":     err.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message": "Pipe has been added to your collection successfully",
		})
		return

	default:
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message": "Operation not allowed for pipe sharing",
		})
		return
	}

}

func (h pipeShareHandler) RemoveShareAccessFromPipe(c *gin.Context) {}
func (h pipeShareHandler) ChangePipeShareAccessType(c *gin.Context) {}
