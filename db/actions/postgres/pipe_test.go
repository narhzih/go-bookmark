package postgres

import (
	"gotest.tools/assert"
	"testing"
)

func Test_pipe_PipeAlreadyExists(t *testing.T) {
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

func Test_pipe_CreatePipe(t *testing.T) {
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

func Test_pipe_GetPipe(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := getPipeTestCases
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			pa := NewPipeActions(db, logger)
			gotPipe, gotErr := pa.GetPipe(tc.inputPipeId, tc.inputUserId)
			assert.Equal(t, tc.wantErr, gotErr)

			if nil == gotErr {
				assert.Equal(t, tc.wantPipe.ID, gotPipe.ID)
				assert.Equal(t, tc.wantPipe.UserID, gotPipe.UserID)
				assert.Equal(t, tc.wantPipe.Name, gotPipe.Name)
			}
		})
	}
}

func Test_pipe_GetPipeByName(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := getPipeByNameTestCases
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			pa := NewPipeActions(db, logger)
			gotPipe, gotErr := pa.GetPipeByName(tc.inputPipeName, tc.inputUserId)
			assert.Equal(t, tc.wantErr, gotErr)

			if nil == gotErr {
				assert.Equal(t, tc.wantPipe.ID, gotPipe.ID)
				assert.Equal(t, tc.wantPipe.UserID, gotPipe.UserID)
				assert.Equal(t, tc.wantPipe.Name, gotPipe.Name)
				assert.Equal(t, tc.wantPipe.Creator, gotPipe.Creator)
			}
		})
	}
}

func Test_pipe_GetPipeAndResource(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := getPipeAndResourceTestCases
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			pa := NewPipeActions(db, logger)
			gotResult, gotErr := pa.GetPipeAndResource(tc.inputPipeId, tc.inputUserId)
			assert.Equal(t, tc.wantErr, gotErr)

			if nil == gotErr {
				assert.Equal(t, tc.wantResult.Pipe.ID, gotResult.Pipe.ID)
				assert.Equal(t, tc.wantResult.Pipe.Name, gotResult.Pipe.Name)
				assert.Equal(t, len(tc.wantResult.Bookmarks), len(gotResult.Bookmarks))
			}
		})
	}
}

func Test_pipe_GetPipes(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := getPipesTestCases
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			pa := NewPipeActions(db, logger)
			gotPipes, gotErr := pa.GetPipes(tc.inputUserId)
			assert.Equal(t, tc.wantErr, gotErr)

			if nil == gotErr {
				assert.Equal(t, len(tc.wantPipes), len(gotPipes))
				assert.Equal(t, tc.wantPipes[0].Name, gotPipes[0].Name)
				assert.Equal(t, tc.wantPipes[0].Creator, gotPipes[0].Creator)
			}
		})
	}
}
