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
		/*
			Algorithm used for public pipe share logic
			---------------------------------------------------------------------
			Check if the user has previously shared this pipe publicly
			---| If they have, return the code for the previous pipe share record
			---| If they haven't, create another record for a public pipe share record and return the code
		*/
		var publicPipeShareRecord models.SharedPipe
		publicPipeShareRecord, err := h.app.Repositories.PipeShare.GetSharedPipe(pipeId, models.PipeShareTypePublic)
		if err != nil {
			if err == postgres.ErrNoRecord {
				// This means no public pipe share record was found for this pipe
				// We can proceed to create a new public pipe share record at this point
				publicPipeShareRecord, err = h.app.Services.SharePipePublicly(pipeId, c.GetInt64(middlewares.KeyUserId))
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"message": "Our system encountered an error while trying to create a public share link",
						"err":     err.Error(),
					})
					return
				}
			} else {
				// An operational error has probably occurred at this point
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Our system encountered an error while trying to create a public share link",
					"err":     err.Error(),
				})
				return
			}
		}

		// Regardless, at this point, it means we have a valid public pipe share code to
		// return to whoever is making this request
		c.JSON(http.StatusCreated, gin.H{
			"message": "Public pipe share link generated successfully",
			"data": map[string]interface{}{
				"share_code": publicPipeShareRecord.Code,
			},
		})
		return

	case models.PipeShareTypePrivate:
		/*
			Algorithm used for private pipe share logic
			---------------------------------------------------------------------
			Check if this pipe has been shared publicly before now
			---| If it has, check if the user has added this pipe to their collection through that public share link
			--- ---| If it has just return a message  telling the user about the situation
			--- ---| If not, Just proceed to the next check


			Check if this pipe has already been shared with the designated req.Username privately prior to now
			---| If it has, return an error message
			---| If it hasn't, create a shared_pipe_receivers record for the designated req.Username
		*/
		if strings.TrimSpace(req.Username) == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "For private pipe share, please specify with whom you want to share the pipe with",
			})
			return
		}
		receiver, err := h.app.Repositories.User.GetUserByUsername(req.Username)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "User not found",
			})
			return
		}
		var pipeHasPublicShareRecord bool
		publicPipeShareRecord, err := h.app.Repositories.PipeShare.GetSharedPipe(pipeId, models.PipeShareTypePublic)
		if err != nil {
			if err == postgres.ErrNoRecord {
				pipeHasPublicShareRecord = true
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Our system encountered an error while trying to create a private share",
					"err":     err.Error(),
				})
				return
			}
		}

		// If this pipe has been shared publicly before now,
		if pipeHasPublicShareRecord {

			// check if the designated req.Username has added this pipe to their collection through that public share link
			var userHasPreviouslyReceivedPipeThroughPublicShareLink = true
			_, err := h.app.Repositories.PipeShare.GetReceivedPipeRecordByCode(publicPipeShareRecord.Code, receiver.ID)
			if err != nil {
				if err == postgres.ErrNoRecord {
					userHasPreviouslyReceivedPipeThroughPublicShareLink = false
				} else {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"message": "Our system encountered an error while trying to create a private share",
						"err":     err.Error(),
					})
				}
			}
			if userHasPreviouslyReceivedPipeThroughPublicShareLink {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": "This pipe cannot be shared with " + receiver.Username + " because they already added this pipe to their collection through a public share link you created",
				})
				return
			}
		}

		// Check if this pipe has already been shared with the designated req.Username privately prior to now
		_, err = h.app.Repositories.PipeShare.GetReceivedPipeRecord(pipeId, receiver.ID)
		if err != nil {
			if err == postgres.ErrNoRecord {
				// This means this pipe has not been shared privately with this user before now
				// Go ahead to create a new private share record
				newPrivatePipeShareRecord := models.SharedPipe{
					PipeID:   pipeId,
					SharerID: c.GetInt64(middlewares.KeyUserId),
				}
				newPrivatePipeShareRecord, err = h.app.Services.SharePipePrivately(newPrivatePipeShareRecord, receiver.Username)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"message": "Our system encountered an error while trying to create a private share",
						"err":     err.Error(),
					})
				}

				err = h.app.Services.CreatePrivatePipeShareNotification(newPrivatePipeShareRecord.Code, newPrivatePipeShareRecord.PipeID, newPrivatePipeShareRecord.SharerID, receiver.ID)
				c.JSON(http.StatusOK, gin.H{
					"message": "Pipe has been successfully shared with " + receiver.Username,
				})
				return
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Our system encountered an error while trying to create a private share",
					"err":     err.Error(),
				})
			}
		}

		// At this point, it means this pipe has already been shared with req.Username
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "This pipe has been previously shared with " + receiver.Username,
		})
		return
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

	/*
		Algorithm used to determine if a user can preview a pipe
		---------------------------------------------------------------------
		Anybody can view a publicly shared pipe, however, for a private pipe
		only the user it was shared with initially can view it
	*/
	if pipeToAdd.Type == models.PipeShareTypePrivate {
		_, err = h.app.Repositories.PipeShare.GetReceivedPipeRecordByCode(code, authenticatedUser.ID)
		if err != nil {
			if err == postgres.ErrNoRecord {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "You cannot view this pipe because it's a private pipe",
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Our system encountered an error while trying to add pipe your collection. Try again soon!",
			})
			return
		}

	}

	// At this point, this pipe can be previewed successfully
	sharer, err := h.app.Repositories.User.GetUserById(pipeToAdd.SharerID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Our system encountered an error while trying to add pipe your collection. Try again soon!",
		})
		return
	}
	pipeAndR, err := h.app.Repositories.Pipe.GetPipeAndResource(pipeToAdd.PipeID, pipeToAdd.SharerID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Our system encountered an error while trying to add pipe your collection. Try again soon!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Preview Successful",
		"data": map[string]interface{}{
			"sharer": sharer,
			"fullPipeData": map[string]interface{}{
				"pipe":      pipeAndR.Pipe,
				"bookmarks": pipeAndR.Bookmarks,
			},
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

	/*
		Algorithm used for private pipe share logic
		---------------------------------------------------------------------
		For publicly shared pipes
		---| Check if the user has added this pipe to their collection through a private share
		--- ---| If they have, return an error message stating they can't add to their collection anymore
		--- ---| If they haven't, check if they have not previously added it to their collection through this same public share link
		--- --- ---| If they have, return error as appropriate
		--- --- ---| If they haven't, add it to their collection


		For privately shared pipes
		---| Check if they've added it to their collection already (either through a private or public share link)
		--- ---| If they have, return error as appropriate,
		--- ---| If they haven't, add it to their collection
	*/
	switch pipeToAdd.Type {
	case models.PipeShareTypePublic:
		receiverRecord, err := h.app.Repositories.PipeShare.GetReceivedPipeRecord(pipeToAdd.PipeID, c.GetInt64(middlewares.KeyUserId))
		if err != nil {
			if err == postgres.ErrNoRecord {
				newReceiverRecord := models.SharedPipeReceiver{
					IsAccepted:   true,
					SharedPipeId: pipeToAdd.PipeID,
					Code:         pipeToAdd.Code,
					SharerId:     pipeToAdd.SharerID,
					ReceiverID:   c.GetInt64(middlewares.KeyUserId),
				}
				newReceiverRecord, err = h.app.Repositories.PipeShare.CreatePipeReceiver(newReceiverRecord)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"message": "Our system encountered an error while trying to add pipe your collection. Try again soon!",
						"err":     err.Error(),
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"message": "Pipe has been added to your collection successfully",
				})
				return
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Our system encountered an error while trying to add pipe your collection. Try again soon!",
					"err":     err.Error(),
				})
				return
			}
		}

		if !receiverRecord.IsAccepted {
			receiverRecord, err = h.app.Repositories.PipeShare.AcceptPrivateShare(receiverRecord)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Our system encountered an error while trying to add pipe your collection. Try again soon!",
					"err":     err.Error(),
				})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Pipe has been added to your collection successfully",
		})
		return

	case models.PipeShareTypePrivate:
		// See if there's any receiver record at all in the database
		receiverRecord, err := h.app.Repositories.PipeShare.GetReceivedPipeRecord(pipeToAdd.PipeID, c.GetInt64(middlewares.KeyUserId))
		if err != nil {
			if err == postgres.ErrNoRecord {
				// If we can't find any receiver record, it means this pipe wasn't shared to you in the first place
				// This is because, at the sharing stage of a pipe, a receiver record is always created for the the receiver
				// only that it needs to be accepted
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "Operation not allowed",
				})
				return
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Our system encountered an error while trying to add pipe your collection. Try again soon!",
					"err":     err.Error(),
				})
				return
			}
		}
		if !receiverRecord.IsAccepted {
			receiverRecord, err = h.app.Repositories.PipeShare.AcceptPrivateShare(receiverRecord)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Our system encountered an error while trying to add pipe your collection. Try again soon!",
					"err":     err.Error(),
				})
				return
			}

		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Pipe has been added to your collection successfully",
		})
		return

		// add a privately shared pipe to your collection
	default:
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message": "Operation not allowed for pipe sharing",
		})
		return
	}

}

func (h pipeShareHandler) RemoveShareAccessFromPipe(c *gin.Context) {}
func (h pipeShareHandler) ChangePipeShareAccessType(c *gin.Context) {}
