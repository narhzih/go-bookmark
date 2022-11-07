package postgres

import (
	"github.com/mypipeapp/mypipeapi/db/models"
	"time"
)

var createAccountVerificationTestCases = map[string]struct {
	inputRecord models.AccountVerification
	wantRecord  models.AccountVerification
	wantErr     error
}{
	"success": {
		inputRecord: models.AccountVerification{
			UserID:    1,
			Token:     "Random_Token_16",
			ExpiresAt: time.Now().Add(120 * time.Second),
		},
		wantRecord: models.AccountVerification{
			ID:     4,
			UserID: 1,
			Token:  "Random_Token_16",
		},
		wantErr: nil,
	},
	"duplicate token": {
		inputRecord: models.AccountVerification{
			UserID: 1,
			Token:  "Random_Token_1",
		},
		wantRecord: models.AccountVerification{},
		wantErr:    ErrRecordExists,
	},
}

var getAccountVerificationByTokenTestCases = map[string]struct {
	inputToken string
	wantRecord models.AccountVerification
	wantErr    error
}{
	"success": {
		inputToken: "Random_Token_1",
		wantRecord: models.AccountVerification{
			UserID: 1,
			Token:  "Random_Token_1",
			ID:     1,
		},
		wantErr: nil,
	},
	"invalid token": {
		inputToken: "Random_Token_15",
		wantRecord: models.AccountVerification{},
		wantErr:    ErrNoRecord,
	},
}

var deleteAccountVerificationTestCases = map[string]struct {
	inputToken   string
	wantResponse bool
	wantErr      error
}{
	"success": {
		inputToken:   "Random_Token_2",
		wantResponse: true,
		wantErr:      nil,
	},
}
