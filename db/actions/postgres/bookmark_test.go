package postgres

import (
	"testing"
)

func Test_bookmark_Create(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}
}
