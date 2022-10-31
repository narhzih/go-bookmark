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

func Test_user_GetUserByTwitterID(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := getUserByTwitterIDTestCases
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			ua := NewUserActions(db, logger)
			gotUser, gotErr := ua.GetUserByTwitterID(tc.inputUserTwitterID)
			assert.Equal(t, tc.wantErr, gotErr)

			if nil == gotErr {
				assert.Equal(t, tc.wantUser.ID, gotUser.ID)
				assert.Equal(t, tc.wantUser.TwitterId, gotUser.TwitterId)
				assert.Equal(t, tc.wantUser.Email, gotUser.Email)
				assert.Equal(t, tc.wantUser.Username, gotUser.Username)
			}
		})
	}
}

func Test_user_GetUserById(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := getUserByIdTestCases
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			ua := NewUserActions(db, logger)
			gotUser, gotErr := ua.GetUserById(tc.inputUserId)
			assert.Equal(t, tc.wantErr, gotErr)

			if nil == gotErr {
				assert.Equal(t, tc.wantUser.ID, gotUser.ID)
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

			if nil == gotErr {
				assert.Equal(t, testCase.wantUser.Email, gotUser.Email)
				assert.Equal(t, testCase.wantUser.Username, gotUser.Username)
			}
		})
	}
}

func Test_user_GetUserByUsername(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := getUserByUsernameTestCases
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			userA := NewUserActions(db, logger)
			gotUser, gotErr := userA.GetUserByUsername(testCase.inputUserUsername)
			assert.Equal(t, testCase.wantErr, gotErr)

			if nil == gotErr {
				assert.Equal(t, testCase.wantUser.Email, gotUser.Email)
				assert.Equal(t, testCase.wantUser.Username, gotUser.Username)
			}
		})
	}
}

func Test_user_GetUserAndAuth(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := getUserAndAuthTestCases
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			userA := NewUserActions(db, logger)
			gotUser, gotErr := userA.GetUserAndAuth(testCase.inputUserId)
			assert.Equal(t, testCase.wantErr, gotErr)

			if nil == gotErr {
				assert.Equal(t, testCase.wantUser.HashedPassword, gotUser.HashedPassword)
				assert.Equal(t, testCase.wantUser.User.Email, gotUser.User.Email)
				assert.Equal(t, testCase.wantUser.User.Username, gotUser.User.Username)
			}
		})
	}
}

func Test_user_GetUserDeviceTokens(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := getUserDeviceTokensTestCases
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			userA := NewUserActions(db, logger)
			gotDeviceTokens, gotErr := userA.GetUserDeviceTokens(testCase.inputUserId)
			assert.Equal(t, testCase.wantErr, gotErr)

			if nil == gotErr {
				assert.Equal(t, len(testCase.wantDeviceTokens), len(gotDeviceTokens))
				assert.Equal(t, testCase.wantDeviceTokens[0], gotDeviceTokens[0])
				assert.Equal(t, testCase.wantDeviceTokens[1], gotDeviceTokens[1])
			}
		})
	}
}

func Test_user_UpdateUser(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := updatedUserTestCases
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			userA := NewUserActions(db, logger)
			gotUser, gotErr := userA.UpdateUser(testCase.inputUser)
			assert.Equal(t, testCase.wantErr, gotErr)

			if nil == gotErr {
				assert.Equal(t, testCase.wantUser.Email, gotUser.Email)
				assert.Equal(t, testCase.wantUser.Username, gotUser.Username)
			}
		})
	}
}

func Test_user_UpdateUserDeviceTokens(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := updateUserDeviceTokens
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			userA := NewUserActions(db, logger)
			gotDeviceTokens, gotErr := userA.UpdateUserDeviceTokens(testCase.inputUserId, testCase.inputUserDeviceTokens)
			assert.Equal(t, testCase.wantErr, gotErr)

			if nil == gotErr {
				assert.Equal(t, len(testCase.wantUserDeviceTokens), len(gotDeviceTokens))
				assert.Equal(t, testCase.wantUserDeviceTokens[0], gotDeviceTokens[0])
				assert.Equal(t, testCase.wantUserDeviceTokens[1], gotDeviceTokens[1])
				assert.Equal(t, testCase.wantUserDeviceTokens[2], gotDeviceTokens[2])
			}
		})
	}
}
