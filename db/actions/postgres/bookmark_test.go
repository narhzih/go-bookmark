package postgres

import (
	"fmt"
	"testing"
)

func Test_bookmark_Create(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	testCases := bookmarkTestCases

	for name, tc := range testCases {
		fmt.Println(name)
		fmt.Println(tc)

		// 1. Create test database & logger
		// 2. Create new bookmark action
		// 3. Call functions and test responses
	}

}
