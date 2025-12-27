package gitt

import "fmt"

var (
	ErrChangesExist = fmt.Errorf("local changes exist")
)
