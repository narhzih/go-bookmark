package postgres

import (
	"github.com/mypipeapp/mypipeapi/db/models"
)

var createUserByEmailTestCases = map[string]struct {
	inputUser models.User
	wantUser  models.User
	wantErr   error
}{
	"success": {
		inputUser: models.User{
			Email:       "user5@gmail.com",
			Username:    "user5",
			ProfileName: "user5",
		},
		wantUser: models.User{
			ID:          5,
			Email:       "user5@gmail.com",
			Username:    "user5",
			ProfileName: "user5",
			CovertPhoto: "",
		},
		wantErr: nil,
	},
	"duplicate name": {
		inputUser: models.User{
			Email:       "user5@gmail.com",
			Username:    "user1",
			ProfileName: "user1",
		},
		wantUser: models.User{},
		wantErr:  ErrRecordExists,
	},

	"duplicate email": {
		inputUser: models.User{
			Email:       "user1@gmail.com",
			Username:    "user5",
			ProfileName: "user5",
		},
		wantUser: models.User{},
		wantErr:  ErrRecordExists,
	},
}

var getUserByTwitterIDTestCases = map[string]struct {
	inputUserTwitterID string
	wantUser           models.User
	wantErr            error
}{
	"success": {
		inputUserTwitterID: "1234567890",
		wantUser: models.User{
			ID:        1,
			Email:     "user1@gmail.com",
			Username:  "user1",
			TwitterId: "1234567890",
		},
		wantErr: nil,
	},
	"invalid twitter id": {
		inputUserTwitterID: "12345",
		wantUser:           models.User{},
		wantErr:            ErrNoRecord,
	},
}

var getUserByIdTestCases = map[string]struct {
	inputUserId int64
	wantUser    models.User
	wantErr     error
}{
	"success": {
		inputUserId: 1,
		wantUser: models.User{
			ID:        1,
			Email:     "user1@gmail.com",
			Username:  "user1",
			TwitterId: "1234567890",
		},
		wantErr: nil,
	},
	"invalid id": {
		inputUserId: 10,
		wantUser:    models.User{},
		wantErr:     ErrNoRecord,
	},
}

var getUserByEmailTestCases = map[string]struct {
	inputUserEmail string
	wantUser       models.User
	wantErr        error
}{
	"success": {
		inputUserEmail: "user1@gmail.com",
		wantUser: models.User{
			ID:          5,
			Email:       "user1@gmail.com",
			Username:    "user1",
			ProfileName: "user1",
			CovertPhoto: "",
		},
		wantErr: nil,
	},
	"invalid email": {
		inputUserEmail: "no_user@gmail.com",
		wantUser:       models.User{},
		wantErr:        ErrNoRecord,
	},
}

var getUserByUsernameTestCases = map[string]struct {
	inputUserUsername string
	wantUser          models.User
	wantErr           error
}{
	"success": {
		inputUserUsername: "user1",
		wantUser: models.User{
			ID:          1,
			Email:       "user1@gmail.com",
			Username:    "user1",
			ProfileName: "user1",
			CovertPhoto: "",
		},
		wantErr: nil,
	},
	"invalid username": {
		inputUserUsername: "no_user",
		wantUser:          models.User{},
		wantErr:           ErrNoRecord,
	},
}

var getUserAndAuthTestCases = map[string]struct {
	inputUserId int64
	wantUser    models.UserAuth
	wantErr     error
}{
	"success": {
		inputUserId: 1,
		wantUser: models.UserAuth{
			HashedPassword: "$2y$15$salteadoususuueyryy28u48viMdUKIwgSc.ETLYvODrrv3MFczPq",
			Origin:         "DEFAULT",
			User: models.User{
				ID:       1,
				Username: "user1",
				Email:    "user1@gmail.com",
			},
		},
		wantErr: nil,
	},
	"invalid user_id": {
		inputUserId: 10,
		wantUser:    models.UserAuth{},
		wantErr:     ErrNoRecord,
	},
}

var getUserDeviceTokensTestCases = map[string]struct {
	inputUserId      int64
	wantDeviceTokens []string
	wantErr          error
}{
	"success": {
		inputUserId:      1,
		wantDeviceTokens: []string{"123", "456"},
		wantErr:          nil,
	},
}
