package postgres

import (
	"github.com/rs/zerolog"
	"gotest.tools/assert"
	"os"
	"testing"
)

func Test_bookmark_Create(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := bookmarkTestCases
	logger := zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()

	for name, tc := range testCases {
		db := newTestDb(t)
		t.Run(name, func(t *testing.T) {
			bookmark := NewBookmarkActions(db, logger)
			gotBookmark, gotErr := bookmark.CreateBookmark(tc.inputBookmark)
			// make sure that we got the expected error
			assert.Equal(t, tc.wantErr, gotErr)

			// in the case of a successful creation
			if nil == gotErr {
				assert.Equal(t, tc.wantBookmark, gotBookmark)
			}
		})
	}

}
