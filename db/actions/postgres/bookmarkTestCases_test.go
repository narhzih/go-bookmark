package postgres

import "github.com/mypipeapp/mypipeapi/db/models"

var bookmarkTestCases = map[string]struct {
	inputBookmark models.Bookmark
	wantBookmark  models.Bookmark
	wantErr       error
}{
	"successful bookmark creation": {
		inputBookmark: models.Bookmark{},
		wantBookmark:  models.Bookmark{},
		wantErr:       nil,
	},
	"existing bookmark": {
		inputBookmark: models.Bookmark{},
		wantBookmark:  models.Bookmark{},
		wantErr:       nil,
	},
}
