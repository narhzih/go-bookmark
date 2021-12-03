package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/trencetech/mypipe-api/db"
	"gitlab.com/trencetech/mypipe-api/db/model"
)

func (h *Handler) CreatePipe(c *gin.Context) {
	newPipeReq := struct {
		Name       string `json:"name" binding:"required"`
		CoverPhoto string `json:"cover_photo" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&newPipeReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}
	pipe := model.Pipe{
		UserID:     c.GetInt64(KeyUserId),
		Name:       newPipeReq.Name,
		CoverPhoto: newPipeReq.CoverPhoto,
	}
	newPipe, err := h.service.DB.CreatePipe(pipe)

	// The only error expected to come up when trying to create a pipe
	// is if there's already a pipe with the same name existing for the
	// user that's trying to create the pipe. This error will be handled later
	// TODO: @narhzih - Implement error handling for UNIQUE(user_id, name)
	if err != nil {
		h.logger.Err(err).Msg(err.Error())
		h.logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred while trying to create pipe",
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

func (h *Handler) GetPipeWithResource(c *gin.Context) {}
func (h *Handler) GetPipes(c *gin.Context) {
	userID := c.GetInt64(KeyUserId)
	pipes, err := h.service.DB.GetPipes(userID)
	if err != nil {
		h.logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"mesage": "An error occurred while trying to fetch pipes",
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
func (h *Handler) UpdatePipe(c *gin.Context) {
	updateReq := struct {
		Name       string `json:"name"`
		CoverPhoto string `json:"cover_photo"`
	}{}

	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	var pipe model.Pipe
	pipeId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid pipe ID",
		})
		return
	}
	pipe = model.Pipe{
		Name:       updateReq.Name,
		CoverPhoto: updateReq.CoverPhoto,
	}
	pipe, err = h.service.DB.UpdatePipe(c.GetInt64(KeyUserId), pipeId, pipe)
	if err != nil {
		h.logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred while trying to update pipe",
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
