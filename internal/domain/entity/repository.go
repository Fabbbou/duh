package entity

const (
	UpdateSafe  = "safe"
	UpdateForce = "force"
	UpdateKeep  = "keep"
)

type Repository struct {
	Name    string
	Aliases map[string]string
	Exports map[string]string
}

type RepositoryUpdateResults struct {
	LocalChangesDetected []string
	OtherErrors          []error
}
