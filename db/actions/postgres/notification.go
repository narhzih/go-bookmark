package postgres

import (
	"database/sql"
	"github.com/mypipeapp/mypipeapi/db/models"
	"github.com/mypipeapp/mypipeapi/db/repository"
	"github.com/rs/zerolog"
)

type notificationActions struct {
	Db     *sql.DB
	Logger zerolog.Logger
}

func NewNotificationActions(db *sql.DB, logger zerolog.Logger) repository.NotificationRepository {
	return notificationActions{
		Db:     db,
		Logger: logger,
	}
}

func (n notificationActions) CreateNotification(userId int64, message, metadata string) (models.Notification, error) {
	var notification models.Notification
	query := `INSERT INTO notifications (user_id, message, metadata) VALUES ($1, $2, $3) RETURNING id, user_id, message, read, metadata, created_at`
	err := n.Db.QueryRow(query, userId, message, metadata).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.Message,
		&notification.Read,
		&notification.MetaData,
		&notification.CreatedAt,
	)

	if err != nil {
		return models.Notification{}, err
	}
	return notification, nil
}

func (n notificationActions) GetNotifications(userId int64) ([]models.Notification, error) {
	var notifications []models.Notification
	query := "SELECT id, user_id, message, read, metadata, created_at FROM notifications WHERE user_id=$1"
	rows, err := n.Db.Query(query, userId)
	if err != nil {
		return notifications, err
	}

	defer rows.Close()
	for rows.Next() {
		var notification models.Notification
		if err := rows.Scan(&notification.ID, &notification.UserID, &notification.Message, &notification.Read, &notification.MetaData, &notification.CreatedAt); err != nil {
			return notifications, err
		}
		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		return notifications, err
	}

	return notifications, nil

}

func (n notificationActions) GetNotification(notificationId, userId int64) (models.Notification, error) {
	var notification models.Notification
	query := "SELECT id, user_id, message, read, metadata, created_at FROM notifications WHERE id=$1 AND user_id=$2 LIMIT 1"
	err := n.Db.QueryRow(query, notificationId, userId).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.Message,
		&notification.Read,
		&notification.MetaData,
		&notification.CreatedAt,
	)
	if err != nil {
		return notification, err
	}
	return notification, nil
}

func (n notificationActions) MarkAsRead(notification models.Notification) (models.Notification, error) {
	var markedNotification models.Notification
	query := `UPDATE notifications SET read=true WHERE id=$1 RETURNING id, user_id, message, metadata, read, created_at`
	err := n.Db.QueryRow(query, notification.ID).Scan(
		&markedNotification.ID,
		&markedNotification.UserID,
		&markedNotification.Message,
		&markedNotification.MetaData,
		&markedNotification.Read,
		&markedNotification.CreatedAt,
	)

	if err != nil {
		return markedNotification, err
	}

	return markedNotification, nil
}
