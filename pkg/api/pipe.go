package api

import (
	"fmt"
	"github.com/rs/zerolog"
	"gitlab.com/trencetech/mypipe-api/pkg/service"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/trencetech/mypipe-api/db"
	"gitlab.com/trencetech/mypipe-api/db/model"
)

func (h *Handler) CreatePipe(c *gin.Context) {
	pipeName := c.PostForm("name")
	if len(pipeName) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"message": "Please specify a pipe name",
		})
		return
	}

	pipeAlreadyExists, err := h.service.DB.PipeAlreadyExists(pipeName, c.GetInt64(KeyUserId))
	if err != nil {
		h.logger.Err(err).Msg(err.Error())
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
	uploadInformation := service.FileUploadInformation{
		Logger:        logger,
		Ctx:           c,
		FileInputName: "cover_photo",
		Type:          "pipe",
	}
	photoUrl, err := service.UploadToCloudinary(uploadInformation)
	if err != nil {
		h.logger.Err(err).Msg(err.Error())
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

	pipe := model.Pipe{
		UserID:     c.GetInt64(KeyUserId),
		Name:       pipeName,
		CoverPhoto: photoUrl,
	}
	newPipe, err := h.service.DB.CreatePipe(pipe)

	// The only error expected to come up when trying to create a pipe
	// is if there's already a pipe with the same name existing for the
	// user that's trying to create the pipe. This error will be handled later
	// TODO: @narhzih - Implement error handling for UNIQUE(user_id, name)
	if err != nil {
		h.logger.Err(err).Msg(err.Error())
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
			},
		},
	})

}
func (h *Handler) GetPipe(c *gin.Context) {
	userID := c.GetInt64(KeyUserId)
	pipeId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid pipe ID",
		})
		return
	}
	h.logger.Info().Msg(fmt.Sprintf("User ID is %+v", userID))

	pipe, err := h.service.DB.GetPipe(pipeId, userID)
	if err != nil {
		if err == db.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Pipe not found",
			})
			return
		}
		h.logger.Err(err).Msg(err.Error())
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

func (h *Handler) GetPipeWithResource(c *gin.Context) {
	userID := c.GetInt64(KeyUserId)
	pipes, err := h.service.DB.GetPipesOnSteroid(userID)
	if err != nil {
		h.logger.Err(err).Msg(err.Error())
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
func (h *Handler) GetPipes(c *gin.Context) {
	userID := c.GetInt64(KeyUserId)
	pipes, err := h.service.DB.GetPipes(userID)
	if err != nil {
		h.logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred while trying to fetch pipes",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "pipes fetched successfully",
		"data": map[string]interface{}{
			"pipes": pipes,
		},
	})
}
func (h *Handler) UpdatePipe(c *gin.Context) {

	var pipe model.Pipe
	pipeId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid pipe ID",
		})
		return
	}
	pipeName := c.PostForm("name")
	if len(pipeName) > 0 {
		pipeAlreadyExists, err := h.service.DB.PipeAlreadyExists(pipeName, c.GetInt64(KeyUserId))
		if err != nil {
			h.logger.Err(err).Msg(err.Error())
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
			h.logger.Err(err).Msg(fmt.Sprintf("file err : %s", err.Error()))
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
		uploadInformation := service.FileUploadInformation{
			Logger:        h.logger,
			Ctx:           c,
			FileInputName: "cover_photo",
			Type:          "pipe",
		}
		photoUrl, err := service.UploadToCloudinary(uploadInformation)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "An error occurred when trying to save user image",
				"err":     err.Error(),
			})
			return
		}
		pipe.CoverPhoto = photoUrl
	}

	pipe, err = h.service.DB.UpdatePipe(c.GetInt64(KeyUserId), pipeId, pipe)
	if err != nil {
		h.logger.Err(err).Msg(err.Error())
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
func (h *Handler) DeletePipe(c *gin.Context) {
	pipeId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid pipe ID",
		})
		return
	}
	_, err = h.service.DB.DeletePipe(c.GetInt64(KeyUserId), pipeId)
	if err != nil {
		h.logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred while trying to delete pipe!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pipe deleted successfully",
	})

}
