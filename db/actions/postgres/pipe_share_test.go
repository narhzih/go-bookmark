package postgres

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_pipe_share_CreatePipeShareRecord(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := createPipeShareRecordTestCases
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			psa := NewPipeShareActions(db, logger)
			gotData, gotErr := psa.CreatePipeShareRecord(tc.inputShareData, tc.inputShareReceiver)
			assert.Equal(t, gotErr, tc.wantErr)
			if nil == gotErr {
				assert.Equal(t, gotData.Code, tc.wantData.Code)
				assert.WithinDuration(t, time.Now(), gotData.CreatedAt, 15*time.Second)
			}
		})
	}
}

func Test_pipe_share_CreatePipeReceiver(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := createPipeReceiverTestCases
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			psa := NewPipeShareActions(db, logger)
			gotReceiver, gotErr := psa.CreatePipeReceiver(tc.inputReceiver)
			assert.Equal(t, gotErr, tc.wantErr)
			if nil == gotErr {
				assert.Equal(t, gotReceiver.ReceiverID, tc.wantReceiver.ReceiverID)
				assert.WithinDuration(t, time.Now(), gotReceiver.CreatedAt, 15*time.Second)
			}
		})
	}
}
func Test_pipe_share_GetSharedPipe(t *testing.T)         {}
func Test_pipe_share_GetSharedPipeByCode(t *testing.T)   {}
func Test_pipe_share_GetReceivedPipeRecord(t *testing.T) {}
func Test_pipe_share_AcceptPrivateShare(t *testing.T)    {}
