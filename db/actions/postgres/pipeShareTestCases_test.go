package postgres

import (
	"fmt"
	"github.com/mypipeapp/mypipeapi/db/models"
)

var createPipeShareRecordTestCases = map[string]struct {
	inputShareData     models.SharedPipe
	inputShareReceiver string
	wantData           models.SharedPipe
	wantErr            error
}{
	"success": {
		inputShareData: models.SharedPipe{
			SharerID: 1,
			PipeID:   1,
			Type:     "private",
			Code:     "jjfjji9993",
		},
		inputShareReceiver: "user2",
		wantData: models.SharedPipe{
			SharerID: 1,
			PipeID:   1,
			Type:     "private",
			Code:     "jjfjji9993",
		},
		wantErr: nil,
	},

	"invalid pipe share type": {
		inputShareData: models.SharedPipe{
			SharerID: 1,
			PipeID:   1,
			Type:     "invalid",
			Code:     "jjfjji9993",
		},
		inputShareReceiver: "user2",
		wantData:           models.SharedPipe{},
		wantErr:            fmt.Errorf("invalid pipe type share: invalid"),
	},
}

var createPipeReceiverTestCases = map[string]struct {
	inputReceiver models.SharedPipeReceiver
	wantReceiver  models.SharedPipeReceiver
	wantErr       error
}{
	"success": {
		inputReceiver: models.SharedPipeReceiver{
			ReceiverID:   2,
			SharerId:     1,
			SharedPipeId: 2,
			IsAccepted:   true,
		},
		wantReceiver: models.SharedPipeReceiver{
			ReceiverID:   2,
			SharerId:     1,
			SharedPipeId: 2,
			IsAccepted:   true,
		},
		wantErr: nil,
	},
}

var getSharedPipeTestCases = map[string]struct {
	inputPipeId    int64
	wantSharedPipe models.SharedPipe
	wantErr        error
}{
	"success": {
		inputPipeId: 1,
		wantSharedPipe: models.SharedPipe{
			SharerID: 1,
			Code:     "MG78k9lig67",
			Type:     "private",
		},
		wantErr: nil,
	},

	"invalid pipe id": {
		inputPipeId:    100,
		wantSharedPipe: models.SharedPipe{},
		wantErr:        ErrNoRecord,
	},
}

var getSharedPipeByCodeTestCases = map[string]struct {
	inputShareCode string
	wantSharedPipe models.SharedPipe
	wantErr        error
}{
	"success": {
		inputShareCode: "MG78k9lig67",
		wantSharedPipe: models.SharedPipe{
			SharerID: 1,
			Code:     "MG78k9lig67",
			Type:     "private",
		},
		wantErr: nil,
	},

	"invalid share code": {
		inputShareCode: "MG78k9lig6788",
		wantSharedPipe: models.SharedPipe{},
		wantErr:        ErrNoRecord,
	},
}

var getReceivedPipeRecordTestCases = map[string]struct {
	inputPipeId            int64
	inputUserId            int64
	wantReceivedPipeRecord models.SharedPipeReceiver
	wantErr                error
}{
	"success": {
		inputPipeId: 1,
		inputUserId: 2,
		wantReceivedPipeRecord: models.SharedPipeReceiver{
			SharerId:     1,
			SharedPipeId: 1,
			ReceiverID:   2,
			IsAccepted:   true,
		},
		wantErr: nil,
	},
	"invalid pipe or user id": {
		inputUserId:            100,
		inputPipeId:            1,
		wantReceivedPipeRecord: models.SharedPipeReceiver{},
		wantErr:                ErrNoRecord,
	},
}

var acceptPrivateShareTestCases = map[string]struct {
	inputReceiver models.SharedPipeReceiver
	wantReceiver  models.SharedPipeReceiver
	wantErr       error
}{
	"success": {
		inputReceiver: models.SharedPipeReceiver{
			ID:           2,
			SharedPipeId: 2,
			SharerId:     1,
			ReceiverID:   2,
			IsAccepted:   false,
		},
		wantReceiver: models.SharedPipeReceiver{
			ID:           2,
			SharedPipeId: 2,
			SharerId:     1,
			ReceiverID:   2,
			IsAccepted:   true,
		},
		wantErr: nil,
	},
}
