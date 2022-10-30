package postgres

import (
	"github.com/mypipeapp/mypipeapi/db/models"
)

// createUserByEmailTestCases - test cases for CreateUserByEmail operations
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
}

// getUserByEmailTestCases - test cases for GetUserByEmail operations
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
}
