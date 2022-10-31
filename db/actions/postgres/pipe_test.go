package postgres

import (
	"gotest.tools/assert"
	"testing"
)

func Test_user_PipeAlreadyExists(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := pipeAlreadyExistsTestCases

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			pa := NewPipeActions(db, logger)

			gotResponse, gotErr := pa.PipeAlreadyExists(tc.inputPipeName, tc.inputUserId)
			assert.Equal(t, tc.wantErr, gotErr)
			assert.Equal(t, tc.wantResponse, gotResponse)
		})
	}
}

func Test_user_CreatePipe(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := createPipeTestCases

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			pa := NewPipeActions(db, logger)
			gotPipe, gotErr := pa.CreatePipe(tc.inputPipe)
			assert.Equal(t, tc.wantErr, gotErr)

			if nil == gotErr {
				assert.Equal(t, tc.wantPipe.Name, gotPipe.Name)
				assert.Equal(t, tc.wantPipe.UserID, gotPipe.UserID)
			}
		})
	}
}
