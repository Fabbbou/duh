package errorss

import "fmt"

type InfrastructureError struct {
	Message string
}

func (e *InfrastructureError) Error() string {
	return fmt.Sprintf("infrastructure error: %s", e.Message)
}
