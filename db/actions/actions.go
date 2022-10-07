package actions

import "fmt"

var (
	ErrRecordExists = fmt.Errorf("row with the same value already exits")
	ErrNoRecord     = fmt.Errorf("no matching row was found")
	//ErrNoRowsInResultSet = fmt.Errorf("no rows in result set")
)
