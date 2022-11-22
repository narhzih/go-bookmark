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

func Test_pipe_share_GetSharedPipe(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := getSharedPipeTestCases
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			psa := NewPipeShareActions(db, logger)
			gotSharedPipe, gotErr := psa.GetSharedPipe(tc.inputPipeId)
			assert.Equal(t, gotErr, tc.wantErr)
			if nil == gotErr {
				assert.Equal(t, gotSharedPipe.SharerID, tc.wantSharedPipe.SharerID)
				assert.Equal(t, gotSharedPipe.Code, tc.wantSharedPipe.Code)
			}
		})
	}
}

func Test_pipe_share_GetSharedPipeByCode(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := getSharedPipeByCodeTestCases
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			psa := NewPipeShareActions(db, logger)
			gotSharedPipe, gotErr := psa.GetSharedPipeByCode(tc.inputShareCode)
			assert.Equal(t, gotErr, tc.wantErr)
			if nil == gotErr {
				assert.Equal(t, gotSharedPipe.SharerID, tc.wantSharedPipe.SharerID)
				assert.Equal(t, gotSharedPipe.Code, tc.wantSharedPipe.Code)
			}
		})
	}
}

func Test_pipe_share_GetReceivedPipeRecord(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := getReceivedPipeRecordTestCases
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			psa := NewPipeShareActions(db, logger)
			gotReceiverRecord, gotErr := psa.GetReceivedPipeRecord(tc.inputPipeId, tc.inputUserId)
			assert.Equal(t, gotErr, tc.wantErr)
			if nil == gotErr {
				assert.Equal(t, gotReceiverRecord.SharerId, tc.wantReceivedPipeRecord.SharerId)
				assert.Equal(t, gotReceiverRecord.ReceiverID, tc.wantReceivedPipeRecord.ReceiverID)
			}
		})
	}
}

func Test_pipe_share_AcceptPrivateShare(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := acceptPrivateShareTestCases
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := newTestDb(t)
			psa := NewPipeShareActions(db, logger)
			gotReceiver, gotErr := psa.AcceptPrivateShare(tc.inputReceiver)
			assert.Equal(t, gotErr, tc.wantErr)

			if nil == gotErr {
				assert.WithinDuration(t, time.Now(), gotReceiver.ModifiedAt, 15*time.Second)
				assert.Equal(t, gotReceiver.IsAccepted, tc.wantReceiver.IsAccepted)
				assert.Equal(t, gotReceiver.SharerId, tc.wantReceiver.SharerId)
			}
		})
	}
}
