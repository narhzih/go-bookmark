package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/trencetech/mypipe-api/db"
	"gitlab.com/trencetech/mypipe-api/db/model"
	"net/http"
	"strconv"
)

func (h *Handler) GetNotification(c *gin.Context) {
	var notification model.Notification
	notificationId, err := strconv.ParseInt(c.Param("notificationId"), 10, 64)
	if err != nil {
		h.logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Notification ID",
		})
		return
	}
	notification, err = h.service.DB.GetNotification(notificationId, c.GetInt64(KeyUserId))
	if err != nil {
		if err == db.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Notification not found",
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to retrieve notification",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification fetched successfully",
		"data": map[string]interface{}{
			"id":         notification.ID,
			"message":    notification.Message,
			"read":       notification.Read,
			"created_at": notification.CreatedAt,
		},
	})
}

func (h *Handler) GetNotifications(c *gin.Context) {
	userId := c.GetInt64(KeyUserId)
	notifications, err := h.service.DB.GetNotifications(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Could not retrieve notifications! Please try again soon",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notifications fetched",
		"data": map[string]interface{}{
			"notifications": notifications,
		},
	})
}
