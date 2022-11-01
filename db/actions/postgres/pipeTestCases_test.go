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

var getPipeTestCases = map[string]struct {
	inputPipeId int64
	inputUserId int64
	wantPipe    models.Pipe
	wantErr     error
}{
	"success": {
		inputUserId: 1,
		inputPipeId: 1,
		wantPipe: models.Pipe{
			UserID:    1,
			ID:        1,
			Name:      "Youtube Shorts",
			Bookmarks: 2,
			Creator:   "user1",
		},
		wantErr: nil,
	},
	"invalid pipe id": {
		inputUserId: 1,
		inputPipeId: 10,
		wantPipe:    models.Pipe{},
		wantErr:     ErrNoRecord,
	},
	"invalid user id": {
		inputUserId: 5,
		inputPipeId: 1,
		wantPipe:    models.Pipe{},
		wantErr:     ErrNoRecord,
	},
}

var getPipeByNameTestCases = map[string]struct {
	inputPipeName string
	inputUserId   int64
	wantPipe      models.Pipe
	wantErr       error
}{
	"success": {
		inputPipeName: "Youtube Shorts",
		inputUserId:   1,
		wantPipe: models.Pipe{
			UserID:    1,
			ID:        1,
			Name:      "Youtube Shorts",
			Bookmarks: 2,
			Creator:   "user1",
		},
		wantErr: nil,
	},
	"invalid pipe name": {
		inputUserId:   1,
		inputPipeName: "Youtube Shorts Invalid",
		wantPipe:      models.Pipe{},
		wantErr:       ErrNoRecord,
	},
	"invalid user name": {
		inputUserId:   10,
		inputPipeName: "Youtube Shorts",
		wantPipe:      models.Pipe{},
		wantErr:       ErrNoRecord,
	},
}

var getPipeAndResourceTestCases = map[string]struct {
	inputPipeId int64
	inputUserId int64
	wantResult  models.PipeAndResource
	wantErr     error
}{
	"success": {
		inputUserId: 1,
		inputPipeId: 1,
		wantResult: models.PipeAndResource{
			Pipe: models.Pipe{
				Name:    "Youtube Shorts",
				Creator: "user1",
				UserID:  1,
				ID:      1,
			},
			Bookmarks: []models.Bookmark{
				{Url: "https://youtu.be/Acgk_Jl95es", Platform: "youtube", PipeID: 1, UserID: 1},
			},
		},
		wantErr: nil,
	},
}

var getPipesTestCases = map[string]struct {
	inputUserId int64
	wantPipes   []models.Pipe
	wantErr     error
}{
	"success": {
		inputUserId: 1,
		wantPipes: []models.Pipe{
			{Name: "Youtube Shorts", ID: 1, UserID: 1, Creator: "user1"},
			{Name: "TikTok", ID: 2, UserID: 1, Creator: "user1"},
		},
		wantErr: nil,
	},
}

var getPipesCountTestCases = map[string]struct {
	inputUserID   int64
	wantPipeCount int
	wantErr       error
}{
	"success": {
		inputUserID:   1,
		wantPipeCount: 2,
		wantErr:       nil,
	},
}

var updatedPipeTestCases = map[string]struct {
	inputUserId      int64
	inputPipeId      int64
	inputUpdatedBody models.Pipe
	wantPipe         models.Pipe
	wantErr          error
}{
	"success": {
		inputUserId: 1,
		inputPipeId: 1,
		inputUpdatedBody: models.Pipe{
			Name:       "Youtube Shortss",
			CoverPhoto: "https://images.unsplash.com/photo-1611162616475-46b635cb6868?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8MXx8eW91dHViZSUyMGxvZ298ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60",
		},
		wantPipe: models.Pipe{
			ID:         1,
			UserID:     1,
			Name:       "Youtube Shortss",
			CoverPhoto: "https://images.unsplash.com/photo-1611162616475-46b635cb6868?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8MXx8eW91dHViZSUyMGxvZ298ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60",
		},
		wantErr: nil,
	},
	"invalid user or pipe id": {
		inputUserId: 10,
		inputPipeId: 10,
		wantPipe:    models.Pipe{},
		inputUpdatedBody: models.Pipe{
			Name:       "Youtube Shortss",
			CoverPhoto: "https://images.unsplash.com/photo-1611162616475-46b635cb6868?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8MXx8eW91dHViZSUyMGxvZ298ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60",
		},
		wantErr: ErrNoRecord,
	},

	"duplicated pipe name": {
		inputUserId: 1,
		inputPipeId: 1,
		inputUpdatedBody: models.Pipe{
			Name:       "TikTok",
			CoverPhoto: "https://images.unsplash.com/photo-1611162616475-46b635cb6868?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8MXx8eW91dHViZSUyMGxvZ298ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60",
		},
		wantPipe: models.Pipe{},
		wantErr:  ErrRecordExists,
	},
}

var deletePipeTestCases = map[string]struct {
	inputUserId  int64
	inputPipeId  int64
	wantResponse bool
	wantErr      error
}{
	"success": {
		inputUserId:  1,
		inputPipeId:  1,
		wantResponse: true,
		wantErr:      nil,
	},
}
