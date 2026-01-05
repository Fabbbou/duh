package entity

import "fmt"

var (
	ErrCouldNotGetPath = fmt.Errorf("could not get path")
	ErrFSDbInitFailed  = fmt.Errorf("filesystem database initialization failed")
)
