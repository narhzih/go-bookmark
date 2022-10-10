package handlers

import (
	"fmt"
	"github.com/rs/zerolog"
	"gitlab.com/trencetech/mypipe-api/cmd/api/internal"
	"gitlab.com/trencetech/mypipe-api/cmd/api/services"
	"gitlab.com/trencetech/mypipe-api/db/actions/postgres"
	"gitlab.com/trencetech/mypipe-api/db/models"
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
	GetPipeWithResource(c *gin.Context)
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

	pipeAlreadyExists, err := h.app.Repositories.Pipe.PipeAlreadyExists(pipeName, c.GetInt64(KeyUserId))
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
		UserID:     c.GetInt64(KeyUserId),
		Name:       pipeName,
		CoverPhoto: photoUrl,
	}
	newPipe, err := h.app.Repositories.Pipe.CreatePipe(pipe)

	// The only error expected to come up when trying to create a pipe
	// is if there's already a pipe with the same name existing for the
	// user that's trying to create the pipe. This error will be handled later
	// TODO: @narhzih - Implement error handling for UNIQUE(user_id, name)
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
				"id":          newPipe.ID,
				"name":        newPipe.Name,
				"cover_photo": newPipe.CoverPhoto,
				"user_id":     newPipe.UserID,
			},
		},
	})

}
func (h pipeHandler) GetPipe(c *gin.Context) {
	userID := c.GetInt64(KeyUserId)
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

func (h pipeHandler) GetPipeWithResource(c *gin.Context) {
	userID := c.GetInt64(KeyUserId)
	pipes, err := h.app.Repositories.Pipe.GetPipesOnSteroid(userID)
	if err != nil {
		h.app.Logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred while trying to fetch pipes",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "pipes fetched succesfully",
		"data": map[string]interface{}{
			"pipes": pipes,
		},
	})
}
func (h pipeHandler) GetPipes(c *gin.Context) {
	userID := c.GetInt64(KeyUserId)
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

	var pipe models.Pipe
	pipeId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.app.Logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid pipe ID",
		})
		return
	}
	pipeName := c.PostForm("name")
	if len(pipeName) > 0 {
		pipeAlreadyExists, err := h.app.Repositories.Pipe.PipeAlreadyExists(pipeName, c.GetInt64(KeyUserId))
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

		pipe.Name = pipeName
	}

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

	pipe, err = h.app.Repositories.Pipe.UpdatePipe(c.GetInt64(KeyUserId), pipeId, pipe)
	if err != nil {
		h.app.Logger.Err(err).Msg(err.Error())
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
	_, err = h.app.Repositories.Pipe.DeletePipe(c.GetInt64(KeyUserId), pipeId)
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
