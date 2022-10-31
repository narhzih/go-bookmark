package postgres

import "fmt"

var (
	ErrRecordExists       = fmt.Errorf("row with the same value already exits")
	ErrNoRecord           = fmt.Errorf("no matching row was found")
	ErrDuplicateUsername  = fmt.Errorf("user with username already exits")
	ErrDuplicateEmail     = fmt.Errorf("user with email already exits")
	ErrDuplicateTwitterID = fmt.Errorf("user with twitter_id already exits")
	//ErrNoRowsInResultSet = fmt.Errorf("no rows in result set")
)
