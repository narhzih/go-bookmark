package postgres

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateVerification(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := createAccountVerificationTestCases

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			ava := NewAccountVerificationActions(db, logger)
			gotRecord, gotErr := ava.CreateVerification(tc.inputRecord)
			assert.Equal(t, gotErr, tc.wantErr)

			if nil == gotErr {
				assert.WithinDuration(t, time.Now(), gotRecord.CreatedAt, 15*time.Second)
				assert.Equal(t, gotRecord.Token, tc.wantRecord.Token)
			}
		})
	}
}

func TestGetAccountVerification(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := getAccountVerificationByTokenTestCases

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			ava := NewAccountVerificationActions(db, logger)
			gotRecord, gotErr := ava.GetAccountVerificationByToken(tc.inputToken)
			assert.Equal(t, gotErr, tc.wantErr)

			if nil == gotErr {
				assert.Equal(t, gotRecord.Token, tc.wantRecord.Token)
				assert.Equal(t, gotRecord.ID, tc.wantRecord.ID)
				assert.Equal(t, gotRecord.UserID, tc.wantRecord.UserID)
			}
		})
	}
}

func TestDeleteVerification(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := deleteAccountVerificationTestCases

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			ava := NewAccountVerificationActions(db, logger)
			gotResponse, gotErr := ava.DeleteVerification(tc.inputToken)
			assert.Equal(t, gotErr, tc.wantErr)

			if nil == gotErr {
				assert.Equal(t, gotResponse, tc.wantResponse)
			}
		})
	}
}
