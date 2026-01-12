package errorss

var (
	ErrCouldNotGetPath = &InfrastructureError{Message: "could not get path"}
	ErrFSDbInitFailed  = &InfrastructureError{Message: "filesystem database initialization failed"}
)
