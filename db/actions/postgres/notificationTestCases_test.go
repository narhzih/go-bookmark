package postgres

import "github.com/mypipeapp/mypipeapi/db/models"

var createNotificationTestCases = map[string]struct {
	inputUserId      int64
	inputMessage     string
	inputMetadata    string
	wantNotification models.Notification
	wantErr          error
}{
	"success": {
		inputUserId:   1,
		inputMessage:  "Just testing a bullshit notification",
		inputMetadata: `{"message": "All good", "work": "all good"}`,
		wantNotification: models.Notification{
			ID:      5,
			UserID:  1,
			Message: "Just testing a bullshit notification",
		},
		wantErr: nil,
	},
}

var getNotificationsTestCases = map[string]struct {
	inputUserId       int64
	wantNotifications []models.Notification
	wantErr           error
}{
	"success": {
		inputUserId: 1,
		wantNotifications: []models.Notification{
			{UserID: 1, Message: "First test on the notification"},
			{UserID: 1, Message: "Second test on the notification"},
			{UserID: 1, Message: "Third test on the notification"},
			{UserID: 1, Message: "Fourth test on the notification"},
		},
		wantErr: nil,
	},
}
