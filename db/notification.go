package db

import "gitlab.com/trencetech/mypipe-api/db/model"

func (db Database) CreateNotification(userId int64, message string) (model.Notification, error) {
	var notification model.Notification
	query := `INSERT INTO notifications (user_id, message) VALUES ($1, $2) RETURNING id, user_id, message, read, created_at`
	err := db.Conn.QueryRow(query, userId, message).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.Message,
		&notification.Read,
		&notification.CreatedAt,
	)

	if err != nil {
		return model.Notification{}, err
	}
	return notification, nil
}

func (db Database) GetNotifications(userId int64) ([]model.Notification, error) {
	var notifications []model.Notification
	query := "SELECT id, user_id, message, read, created_at FROM notifications WHERE user_id=$1"
	rows, err := db.Conn.Query(query, userId)
	if err != nil {
		return notifications, err
	}

	defer rows.Close()
	for rows.Next() {
		var notification model.Notification
		if err := rows.Scan(&notification.ID, &notification.UserID, &notification.Message, &notification.Read, &notification.CreatedAt); err != nil {
			return notifications, err
		}
		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		return notifications, err
	}

	return notifications, nil

}

func (db Database) GetNotification(notificationId, userId int64) (model.Notification, error) {
	var notification model.Notification
	query := "SELECT id, user_id, message, read, created_at FROM notifications WHERE id=$1 AND user_id=$2 LIMIT 1"
	err := db.Conn.QueryRow(query, notificationId, userId).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.Message,
		&notification.Read,
		&notification.CreatedAt,
	)
	if err != nil {
		return notification, err
	}
	return notification, nil
}
