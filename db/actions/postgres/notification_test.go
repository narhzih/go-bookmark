package postgres

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_notification_CreateNotification(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := createNotificationTestCases
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			na := NewNotificationActions(db, logger)
			gotNotification, gotErr := na.CreateNotification(tc.inputUserId, tc.inputMessage, tc.inputMetadata)
			assert.Equal(t, tc.wantErr, gotErr)

			if nil == gotErr {
				// make sure the notification was created within 15 seconds
				assert.WithinDuration(t, time.Now(), gotNotification.CreatedAt, 15*time.Second)

				assert.Equal(t, tc.wantNotification.UserID, gotNotification.UserID)
				assert.Equal(t, tc.wantNotification.ID, gotNotification.ID)
			}
		})
	}
}

func Test_notification_GetNotifications(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := getNotificationsTestCases
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			na := NewNotificationActions(db, logger)
			gotNotifications, gotErr := na.GetNotifications(tc.inputUserId)
			assert.Equal(t, tc.wantErr, gotErr)

			if nil == gotErr {
				// confirm that the operation completed within 15 seconds

				assert.Equal(t, len(tc.wantNotifications), len(gotNotifications))

			}
		})
	}
}
