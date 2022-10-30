package postgres

import (
	"gotest.tools/assert"
	"testing"
)

func Test_user_CreateUserByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}
	testCases := createUserByEmailTestCases
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			userA := NewUserActions(db, logger)
			gotUser, gotErr := userA.CreateUserByEmail(tc.inputUser, "", "DEFAULT")
			assert.Equal(t, tc.wantErr, gotErr)

			// in the case of a successful creation
			if nil == gotErr {
				assert.Equal(t, tc.wantUser.Email, gotUser.Email)
				assert.Equal(t, tc.wantUser.Username, gotUser.Username)
			}
		})
	}
}

func Test_user_GetUserByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := getUserByEmailTestCases
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			userA := NewUserActions(db, logger)
			gotUser, gotErr := userA.GetUserByEmail(testCase.inputUserEmail)
			assert.Equal(t, testCase.wantErr, gotErr)
			assert.Equal(t, testCase.wantUser.Email, gotUser.Email)
		})
	}
}
