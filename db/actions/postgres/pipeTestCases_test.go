package postgres

import "github.com/mypipeapp/mypipeapi/db/models"

var pipeAlreadyExistsTestCases = map[string]struct {
	inputUserId   int64
	inputPipeName string
	wantResponse  bool
	wantErr       error
}{
	"pipe exists": {
		inputUserId:   1,
		inputPipeName: "Youtube Shorts",
		wantResponse:  true,
		wantErr:       nil,
	},
	"pipe does not exits": {
		inputUserId:   1,
		inputPipeName: "Youtube",
		wantResponse:  false,
		wantErr:       nil,
	},
}

var createPipeTestCases = map[string]struct {
	inputPipe models.Pipe
	wantPipe  models.Pipe
	wantErr   error
}{
	"success": {
		inputPipe: models.Pipe{
			Name:       "Instagram",
			UserID:     1,
			CoverPhoto: "https://images.unsplash.com/photo-1611162616305-c69b3fa7fbe0?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8NHx8eW91dHViZSUyMGxvZ298ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60",
		},
		wantPipe: models.Pipe{
			Name:       "Instagram",
			UserID:     1,
			CoverPhoto: "https://images.unsplash.com/photo-1611162616305-c69b3fa7fbe0?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8NHx8eW91dHViZSUyMGxvZ298ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60",
		},
		wantErr: nil,
	},
	"duplicate name": {
		inputPipe: models.Pipe{
			Name:   "TikTok",
			UserID: 1,
		},
		wantPipe: models.Pipe{},
		wantErr:  ErrRecordExists,
	},
}
