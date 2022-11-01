package handlers

import (
	"encoding/json"
	"fmt"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/mypipeapp/mypipeapi/cmd/api/helpers"
	"github.com/mypipeapp/mypipeapi/cmd/api/internal"
	"github.com/mypipeapp/mypipeapi/cmd/api/middlewares"
	"github.com/mypipeapp/mypipeapi/cmd/api/services"
	"github.com/mypipeapp/mypipeapi/db/actions/postgres"
	"github.com/mypipeapp/mypipeapi/db/models"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PipeHandler interface {
	CreatePipe(c *gin.Context)
	GetPipe(c *gin.Context)
	UpdatePipe(c *gin.Context)
	DeletePipe(c *gin.Context)
	GetPipes(c *gin.Context)
}

type pipeHandler struct {
	app internal.Application
}

func NewPipeHandler(app internal.Application) PipeHandler {
	return pipeHandler{app: app}
}

func (h pipeHandler) CreatePipe(c *gin.Context) {
	pipeName := c.PostForm("name")
	var photoUrl = ""
	if len(pipeName) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"message": "Please specify a pipe name",
		})
		return
	}

	pipeAlreadyExists, err := h.app.Repositories.Pipe.PipeAlreadyExists(pipeName, c.GetInt64(middlewares.KeyUserId))
	if err != nil {
		h.app.Logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "An error occurred during validation",
			"err":     err.Error(),
		})
		return
	}
	if pipeAlreadyExists == true {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "You already have a pipe with this name",
		})
		return
	}
	logger := zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()
	// Check if a file was added
	_, _, err = c.Request.FormFile("cover_photo")
	if err != nil {
		if err != http.ErrMissingFile {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
	} else {
		uploadInformation := services.FileUploadInformation{
			Logger:        logger,
			Ctx:           c,
			FileInputName: "cover_photo",
			Type:          "pipe",
		}
		photoUrl, err = services.UploadToCloudinary(uploadInformation)
		if err != nil {
			h.app.Logger.Err(err).Msg(err.Error())
			if err == http.ErrMissingFile {
				c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
					"message": "No file was uploaded. Please select a file to upload as your pipe cover",
					"err":     err.Error(),
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "An error occurred when trying to process pipe image",
				"err":     err.Error(),
			})
			return
		}
	}
	h.app.Logger.Info().Msg("Actual pipe creation has started")
	pipe := models.Pipe{
		UserID:     c.GetInt64(middlewares.KeyUserId),
		Name:       pipeName,
		CoverPhoto: photoUrl,
	}
	newPipe, err := h.app.Repositories.Pipe.CreatePipe(pipe)
	if err != nil {
		h.app.Logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred while trying to create pipe",
			"err":     err.Error(),
		})
		return
	}
	// The only error expected to come up when trying to create a pipe
	// is if there's already a pipe with the same name existing for the
	// user that's trying to create the pipe. This error will be handled later
	fetchedPipe, err := h.app.Repositories.Pipe.GetPipe(newPipe.ID, newPipe.UserID)
	if err != nil {
		h.app.Logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred while trying to create pipe",
			"err":     err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pipe created successfully",
		"data": map[string]interface{}{
			"pipe": map[string]interface{}{
				"id":          fetchedPipe.ID,
				"user_id":     fetchedPipe.UserID,
				"bookmarks":   fetchedPipe.Bookmarks,
				"name":        fetchedPipe.Name,
				"created_at":  fetchedPipe.CreatedAt,
				"modified_at": fetchedPipe.ModifiedAt,
				"creator":     fetchedPipe.Creator,
			},
		},
	})

}
func (h pipeHandler) GetPipe(c *gin.Context) {
	userID := c.GetInt64(middlewares.KeyUserId)
	pipeId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.app.Logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid pipe ID",
		})
		return
	}
	h.app.Logger.Info().Msg(fmt.Sprintf("User ID is %+v", userID))

	pipe, err := h.app.Repositories.Pipe.GetPipe(pipeId, userID)
	if err != nil {
		if err == postgres.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Pipe not found",
			})
			return
		}
		h.app.Logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "An error occurred while trying to retrieve pipe. Please try again later",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "pipe fetched successfully",
		"data": map[string]interface{}{
			"pipe": map[string]interface{}{
				"id":          pipe.ID,
				"name":        pipe.Name,
				"cover_photo": pipe.CoverPhoto,
			},
		},
	})
}

func (h pipeHandler) GetPipes(c *gin.Context) {
	userID := c.GetInt64(middlewares.KeyUserId)
	pipes, err := h.app.Repositories.Pipe.GetPipes(userID)
	if err != nil {
		h.app.Logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred while trying to fetch pipes",
		})
		return
	}
	if pipes == nil || len(pipes) <= 0 {

	}
	c.JSON(http.StatusOK, gin.H{
		"message": "pipes fetched successfully",
		"data": map[string]interface{}{
			"pipes": pipes,
		},
	})
}
func (h pipeHandler) UpdatePipe(c *gin.Context) {
	req := struct {
		Name string `form:"name" json:"name,omitempty"`
	}{}

	if err := c.ShouldBind(&req); err != nil {
		errMessage := helpers.ParseErrorMessage(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": errMessage,
			"err":     err.Error(),
		})
		return
	}

	var pipe models.Pipe

	pipeId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.app.Logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid pipe ID",
		})
		return
	}

	pipe, err = h.app.Repositories.Pipe.GetPipe(pipeId, c.GetInt64(middlewares.KeyUserId))
	if err != nil {
		if err == postgres.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "invalid pipe",
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "our server encountered an error",
			"err":     err.Error(),
		})
		return
	}

	// creat pipe
	pipeBytes, _ := json.Marshal(&pipe)
	reqBytes, _ := json.Marshal(&req)
	updatedBodyBytes, err := jsonpatch.MergePatch(pipeBytes, reqBytes)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "our server encountered an error",
			"err":     err.Error(),
		})
		return
	}

	_ = json.Unmarshal(updatedBodyBytes, &pipe)

	file, _, err := c.Request.FormFile("cover_photo")
	if err != nil {
		if err != http.ErrMissingFile {
			h.app.Logger.Err(err).Msg(fmt.Sprintf("file err : %s", err.Error()))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "An error occurred",
				"err":     err.Error(),
			})
			return
		}
	}
	if file != nil {
		// This means a file was uploaded with the request
		// Try uploading it to Cloudinary
		uploadInformation := services.FileUploadInformation{
			Logger:        h.app.Logger,
			Ctx:           c,
			FileInputName: "cover_photo",
			Type:          "pipe",
		}
		photoUrl, err := services.UploadToCloudinary(uploadInformation)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "An error occurred when trying to save user image",
				"err":     err.Error(),
			})
			return
		}
		pipe.CoverPhoto = photoUrl
	}

	pipe, err = h.app.Repositories.Pipe.UpdatePipe(c.GetInt64(middlewares.KeyUserId), pipeId, pipe)
	if err != nil {
		if err == postgres.ErrRecordExists {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "You already have a pipe with this name",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred while trying to update pipe",
			"err":     err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pipe updated successfully",
		"data": map[string]interface{}{
			"pipe": map[string]interface{}{
				"id":          pipe.ID,
				"name":        pipe.Name,
				"cover_photo": pipe.CoverPhoto,
			},
		},
	})

}
func (h pipeHandler) DeletePipe(c *gin.Context) {
	pipeId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.app.Logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid pipe ID",
		})
		return
	}
	_, err = h.app.Repositories.Pipe.DeletePipe(c.GetInt64(middlewares.KeyUserId), pipeId)
	if err != nil {
		h.app.Logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred while trying to delete pipe!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pipe deleted successfully",
	})

}
