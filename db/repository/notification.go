package repository

import "gitlab.com/trencetech/mypipe-api/db/models"

type NotificationRepository interface {
	CreateNotification(userId int64, message, metadata string) (models.Notification, error)
	GetNotifications(userId int64) ([]models.Notification, error)
	GetNotification(notificationId, userId int64) (models.Notification, error)
	MarkAsRead(notification models.Notification) (models.Notification, error)
}
