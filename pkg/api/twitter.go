package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/trencetech/mypipe-api/db"
	"gitlab.com/trencetech/mypipe-api/db/model"
	"gitlab.com/trencetech/mypipe-api/pkg/helpers"
	"net/http"
)

func (h *Handler) BotAddToPipe(c *gin.Context) {
	var pipe model.Pipe
	botAddToPipeBody := struct {
		TwitterID string `json:"twitter_id" binding:"required"`
		TweetLink string `json:"tweet_link" binding:"required"`
		PipeName  string `json:"pipe_name" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&botAddToPipeBody); err != nil {
		errMessage := helpers.ParseErrorMessage(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": errMessage,
			"err":     err.Error(),
		})
		return
	}
	user, err := h.service.TwitterAccountConnected(botAddToPipeBody.TwitterID)

	if err != nil {
		if err == db.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"message": `OOPS! Something went wrong! Either you haven't connected your twitter account to your MyPipe account or, you do not have a MyPipe account at all. Either ways, follow this like so I can quickly fix this for you!`,
				"data": map[string]interface{}{
					"link": c.Request.Host + "/v1/bot/twitter/connect-account",
				},
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No connected twitter account",
		})
		return
	}
	pipe, err = h.service.DB.GetPipeByName(botAddToPipeBody.PipeName, user.ID)
	if err != nil {
		if err == db.ErrNoRecord {
			// Create a pipe for that user
			pipe, err = h.service.DB.CreatePipe(model.Pipe{
				Name:       botAddToPipeBody.PipeName,
				UserID:     user.ID,
				CoverPhoto: "",
			})
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Our server encountered an error. Please try again later",
				})
			}
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "An error occurred while trying to perform operation",
				"err":     err.Error(),
			})
			return
		}
	}
	bookmark, err := h.service.DB.CreateBookmark(model.Bookmark{
		UserID:   user.ID,
		PipeID:   pipe.ID,
		Platform: "twitter",
		Url:      botAddToPipeBody.TweetLink,
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Our server encountered an error. Please try again later",
		})
	}

	err = h.service.CreateTwitterPipeShareNotification(bookmark.Url, pipe.Name, user.ID)
	if err != nil {
		h.logger.Err(err).Msg("An error occurred while creating notification for twitter pipe share")
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bookmark created successfully",
		"data": map[string]interface{}{
			"bookmark": bookmark,
		},
	})

}
