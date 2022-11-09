package postgres

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_bookmark_Create(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := createBookmarkTestCases
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			ba := NewBookmarkActions(db, logger)
			gotBookmark, gotErr := ba.CreateBookmark(tc.inputBookmark)
			assert.Equal(t, gotErr, tc.wantErr)

			if gotErr == nil {
				// check if the bookmark was created less than 15 seconds ago
				assert.WithinDuration(t, time.Now(), gotBookmark.CreatedAt, 15*time.Second)

				assert.Equal(t, gotBookmark.ID, tc.wantBookmark.ID)
				assert.Equal(t, gotBookmark.PipeID, tc.wantBookmark.PipeID)
			}
		})
	}
}

func Test_bookmark_GetBookmark(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := getBookmarkTestCases
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			ba := NewBookmarkActions(db, logger)
			gotBookmark, gotErr := ba.GetBookmark(tc.inputBookmarkId, tc.inputUserId)
			assert.Equal(t, gotErr, tc.wantErr)

			if nil == gotErr {
				assert.Equal(t, gotBookmark.ID, tc.wantBookmark.ID)
				assert.Equal(t, gotBookmark.UserID, tc.wantBookmark.UserID)
			}
		})
	}
}

func Test_bookmark_GetBookmarks(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := getBookmarksTestCases
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			ba := NewBookmarkActions(db, logger)
			gotBookmarks, gotErr := ba.GetBookmarks(tc.inputUserId, tc.inputPipeId)
			assert.Equal(t, tc.wantErr, gotErr)

			if nil == gotErr {
				assert.Equal(t, len(tc.wantBookmarks), len(gotBookmarks))
				// Validate the tags for any of the bookmarks as well
			}
		})
	}
}

func Test_bookmark_GetBookmarksCount(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := getBookmarksCountTestCases
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			ba := NewBookmarkActions(db, logger)
			gotBookmarks, gotErr := ba.GetBookmarksCount(tc.inputUserId)
			assert.Equal(t, gotErr, tc.wantErr)

			if nil == gotErr {
				assert.Equal(t, gotBookmarks, tc.wantBookmarks)
			}
		})
	}
}

func Test_bookmark_DeleteBookmark(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := deleteBookmarkTestCases
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			ba := NewBookmarkActions(db, logger)
			gotResponse, gotErr := ba.DeleteBookmark(tc.inputBookmarkId, tc.inputUserId)
			assert.Equal(t, tc.wantErr, gotErr)
			assert.Equal(t, tc.wantResponse, gotResponse)
		})
	}
}
