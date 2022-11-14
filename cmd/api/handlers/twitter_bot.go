package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mypipeapp/mypipeapi/cmd/api/helpers"
	"github.com/mypipeapp/mypipeapi/cmd/api/internal"
	"github.com/mypipeapp/mypipeapi/db/actions/postgres"
	"github.com/mypipeapp/mypipeapi/db/models"
	"net/http"
	"strings"
)

type TwitterBotHandler interface {
	BotAddToPipe(c *gin.Context)
}

type twitterBotHandler struct {
	app internal.Application
}

func NewTwitterBotHandler(app internal.Application) TwitterBotHandler {
	return twitterBotHandler{app: app}
}

func (h twitterBotHandler) BotAddToPipe(c *gin.Context) {
	var pipe models.Pipe
	req := struct {
		TwitterID string `json:"twitter_id" binding:"required"`
		TweetLink string `json:"tweet_link" binding:"required"`
		PipeName  string `json:"pipe_name" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		errMessage := helpers.ParseErrorMessage(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": errMessage,
			"err":     err.Error(),
		})
		return
	}
	user, err := h.app.Services.TwitterAccountConnected(req.TwitterID)

	if err != nil {
		if err == postgres.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"message": `OOPS! Something went wrong! Either you haven't connected your twitter account to your MyPipe account or, you do not have a MyPipe account at all. Either ways, follow this link so I can quickly fix this for you!`,
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
	refinedPipeNames := strings.TrimSpace(req.PipeName)
	pipeNames := strings.Split(refinedPipeNames, ",")
	for _, name := range pipeNames {
		h.app.Logger.Info().Msg(fmt.Sprintf("adding bookmark to %s", name))
		pipe, err = h.app.Repositories.Pipe.GetPipeByName(strings.ToLower(name), user.ID)
		if err != nil {
			if err == postgres.ErrNoRecord {
				h.app.Logger.Info().Msg(fmt.Sprintf("pipe name %s doesn't exist, moving on...", name))
				// Just skip the process for this particular pipe if it doesn't exist
				continue
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "An error occurred while trying to perform operation",
					"err":     err.Error(),
				})
				return
			}
		}
		_, err = h.app.Repositories.Bookmark.CreateBookmark(models.Bookmark{
			UserID:   user.ID,
			PipeID:   pipe.ID,
			Platform: "twitter",
			Url:      req.TweetLink,
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Our server encountered an error. Please try again later",
			})
			return
		}
	}

	//if bookmark.Url != "" {
	//	err = h.app.Services.CreateTwitterPipeShareNotification(bookmark.Url, pipe.Name, user.ID)
	//	if err != nil {
	//		h.app.Logger.Err(err).Msg("An error occurred while creating notification for twitter pipe share")
	//	}
	//}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bookmark created added to desired pipes successfully",
	})
}
