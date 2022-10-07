package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/trencetech/mypipe-api/db"
	"gitlab.com/trencetech/mypipe-api/db/models"
	"net/http"
	"strconv"
)

func (h *Handler) GetNotification(c *gin.Context) {
	var notification models.Notification
	notificationId, err := strconv.ParseInt(c.Param("notificationId"), 10, 64)
	if err != nil {
		h.logger.Err(err).Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Notification ID",
		})
		return
	}
	h.logger.Info().Msg(fmt.Sprintf("retrieving notification for %+v", notificationId))
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
			"err":     err.Error(),
		})
		return
	}

	if !notification.Read {
		// Mark notification as read
		notification, err = h.service.DB.MarkAsRead(notification)
		if err != nil {
			h.logger.Err(err).Msg("an error occurred while trying to mark notification as read")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification fetched successfully",
		"data": map[string]interface{}{
			"id":         notification.ID,
			"message":    notification.Message,
			"metadata":   notification.MetaData,
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

func (h *Handler) UpdateUserDeviceTokens(c *gin.Context) {
	reqBody := struct {
		DeviceToken string `json:"device_token"`
	}{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Please provide a device token",
		})
		return
	}

	existingDeviceTokens, err := h.service.DB.GetUserDeviceTokens(c.GetInt64(KeyUserId))
	if err != nil {
		if err != db.ErrNoRecord {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "An error occurred",
				"error":   err.Error(),
			})
			return
		}
	}
	// TODO: Refactor to remove old device tokens if regenerated
	existingDeviceTokens = append(existingDeviceTokens, reqBody.DeviceToken)
	existingDeviceTokens, err = h.service.DB.UpdateUserDeviceTokens(c.GetInt64(KeyUserId), existingDeviceTokens)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Our system encountered an error. Please try again soon",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User device token updated successfully",
	})

}
