package postgres

import "github.com/mypipeapp/mypipeapi/db/models"

var bookmarkTestCases = map[string]struct {
	inputBookmark models.Bookmark
	gotBookmark   models.Bookmark
	gotErr        error
}{
	"valid bookmark": {
		inputBookmark: models.Bookmark{},
		gotBookmark:   models.Bookmark{},
		gotErr:        nil,
	},
	"invalid bookmark": {
		inputBookmark: models.Bookmark{},
		gotBookmark:   models.Bookmark{},
		gotErr:        nil,
	},
}
