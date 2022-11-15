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
