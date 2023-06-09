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
			HashedPassword: "$2a$14$A/CXTnm0.WSb0CoWcH31VeKv.CitRdGTiWHj/06I3cUvwgrj.UwBu",
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

var updatedUserTestCases = map[string]struct {
	inputUser models.User
	wantUser  models.User
	wantErr   error
}{
	"success": {
		inputUser: models.User{
			ID:          1,
			Username:    "user1_updated",
			Email:       "user1@gmail.com",
			TwitterId:   "1234567890",
			CovertPhoto: "https://images.unsplash.com/photo-1611608822650-925c227ef4d2?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8Nnx8aGFuZHNvbWUlMjBtYW58ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60",
		},
		wantUser: models.User{
			ID:          1,
			Username:    "user1_updated",
			Email:       "user1@gmail.com",
			TwitterId:   "1234567890",
			CovertPhoto: "https://images.unsplash.com/photo-1611608822650-925c227ef4d2?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8Nnx8aGFuZHNvbWUlMjBtYW58ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60",
		},
		wantErr: nil,
	},
	"duplicate username": {
		inputUser: models.User{
			ID:          1,
			Username:    "user2",
			Email:       "user1@gmail.com",
			TwitterId:   "1234567890",
			CovertPhoto: "https://images.unsplash.com/photo-1611608822650-925c227ef4d2?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8Nnx8aGFuZHNvbWUlMjBtYW58ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60",
		},
		wantUser: models.User{},
		wantErr:  ErrDuplicateUsername,
	},
	"duplicate email": {
		inputUser: models.User{
			ID:          1,
			Username:    "user1",
			Email:       "user2@gmail.com",
			TwitterId:   "1234567890",
			CovertPhoto: "https://images.unsplash.com/photo-1611608822650-925c227ef4d2?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8Nnx8aGFuZHNvbWUlMjBtYW58ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60",
		},
		wantUser: models.User{},
		wantErr:  ErrDuplicateEmail,
	},
}

var updateUserDeviceTokens = map[string]struct {
	inputUserId           int64
	inputUserDeviceTokens []string
	wantUserDeviceTokens  []string
	wantErr               error
}{
	"success": {
		inputUserId:           1,
		inputUserDeviceTokens: []string{"123", "456", "789"},
		wantUserDeviceTokens:  []string{"123", "456", "789"},
		wantErr:               nil,
	},
}
