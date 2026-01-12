package entity

const (
	UpdateSafe  = "safe"
	UpdateForce = "force"
	UpdateKeep  = "keep"
)

type Package struct {
	Name                 string
	Aliases              map[string]string
	Exports              map[string]string
	GitConfigIncludePath string
}

type PackageUpdateResults struct {
	LocalChangesDetected []string
	OtherErrors          []error
}
