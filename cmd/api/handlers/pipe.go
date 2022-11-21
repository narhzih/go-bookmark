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
	"net/http"
	"strconv"
	"strings"

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
	req := struct {
		Name string `form:"name" json:"name" binding:"required"`
	}{}

	if err := c.ShouldBind(&req); err != nil {
		errMessage := helpers.ParseErrorMessage(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": errMessage,
		})
	}
	authenticatedUser := middlewares.GetLoggedInUser(c)

	// TODO: refine db call to remove extra step of first checking if pipe
	//		 already exists
	pipeAlreadyExists, err := h.app.Repositories.Pipe.PipeAlreadyExists(req.Name, authenticatedUser.ID)
	if err != nil {
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

	// Check if a file was added
	var photoUrl = ""
	file, _, err := c.Request.FormFile("cover_photo")
	if err != nil {
		if err != http.ErrMissingFile {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
	}

	if file != nil {
		// Upload image to cloudinary
		uploadInformation := services.FileUploadInformation{
			Logger:        h.app.Logger,
			Ctx:           c,
			FileInputName: "cover_photo",
			Type:          "pipe",
		}
		photoUrl, err = services.UploadToCloudinary(uploadInformation)
		if err != nil {
			h.app.Logger.Err(err).Msg(err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "An error occurred when trying to process pipe image",
				"err":     err.Error(),
			})
			return
		}
	}

	pipe := models.Pipe{
		UserID:     authenticatedUser.ID,
		Name:       strings.TrimSpace(strings.ToLower(req.Name)),
		CoverPhoto: photoUrl,
	}

	pipe, err = h.app.Repositories.Pipe.CreatePipe(pipe)
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
	pipe, err = h.app.Repositories.Pipe.GetPipe(pipe.ID, pipe.UserID)
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
			"pipe": pipe,
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
			"pipe": pipe,
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

	pipe.Name = strings.TrimSpace(strings.ToLower(pipe.Name))
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

	pipe, _ = h.app.Repositories.Pipe.GetPipe(pipe.ID, pipe.UserID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pipe updated successfully",
		"data": map[string]interface{}{
			"pipe": pipe,
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
