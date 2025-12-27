package gitt

import "fmt"

const (
	GitStrategySafe   = "safe"
	GitStrategyCommit = "commit"
	GitStrategyForce  = "force"
)

var (
	ErrChangesExist = fmt.Errorf("local changes exist")
)
