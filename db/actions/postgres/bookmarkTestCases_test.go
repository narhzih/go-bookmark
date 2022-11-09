package postgres

import "github.com/mypipeapp/mypipeapi/db/models"

var createBookmarkTestCases = map[string]struct {
	inputBookmark models.Bookmark
	wantBookmark  models.Bookmark
	wantErr       error
}{
	"success": {
		inputBookmark: models.Bookmark{
			UserID:   1,
			Platform: "twitter",
			PipeID:   1,
			Url:      "https://twitter.com/Mc_Phils/status/1589501899015090178?s=20&t=AXekm5YnalcWausr3fuqlA",
		},
		wantBookmark: models.Bookmark{
			ID:       7,
			UserID:   1,
			Platform: "twitter",
			PipeID:   1,
			Url:      "https://twitter.com/Mc_Phils/status/1589501899015090178?s=20&t=AXekm5YnalcWausr3fuqlA",
		},
		wantErr: nil,
	},
}

var getBookmarkTestCases = map[string]struct {
	inputBookmarkId int64
	inputUserId     int64
	wantBookmark    models.Bookmark
	wantErr         error
}{
	"success": {
		inputBookmarkId: 1,
		inputUserId:     1,
		wantBookmark: models.Bookmark{
			ID:       1,
			UserID:   1,
			Platform: "youtube",
			Url:      "https://youtu.be/Acgk_Jl95es",
		},
		wantErr: nil,
	},

	"invalid id": {
		inputBookmarkId: 1000,
		inputUserId:     1,
		wantBookmark:    models.Bookmark{},
		wantErr:         ErrNoRecord,
	},
}

var getBookmarksTestCases = map[string]struct {
	inputPipeId   int64
	inputUserId   int64
	wantBookmarks []models.Bookmark
	wantErr       error
}{
	"success": {
		inputUserId: 1,
		inputPipeId: 1,
		wantBookmarks: []models.Bookmark{
			{
				ID:       1,
				UserID:   1,
				PipeID:   1,
				Url:      "https://youtu.be/Acgk_Jl95es",
				Platform: "youtube",
				Tags:     []string{"Beautiful Asian Muslim", "Quick Blows"},
			},
		},
		wantErr: nil,
	},
	"invalid user or pipe id": {
		inputUserId:   1,
		inputPipeId:   2000,
		wantBookmarks: []models.Bookmark{},
		wantErr:       nil,
	},
}

var getBookmarksCountTestCases = map[string]struct {
	inputUserId   int64
	wantBookmarks int
	wantErr       error
}{
	"success": {
		inputUserId:   1,
		wantBookmarks: 2,
		wantErr:       nil,
	},
}

var deleteBookmarkTestCases = map[string]struct {
	inputBookmarkId int64
	inputUserId     int64
	wantResponse    bool
	wantErr         error
}{
	"success": {
		inputBookmarkId: 1,
		inputUserId:     1,
		wantResponse:    true,
		wantErr:         nil,
	},
}
