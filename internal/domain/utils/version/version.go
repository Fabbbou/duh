package version

import (
	"fmt"
	"runtime"
)

// These variables will be set at build time using ldflags
var (
	// Version is the current version of the application
	Version = "dev"
	Built   = "unknown"
	Commit  = "none"
)

// BuildInfo returns detailed build information
func BuildInfo() string {
	return fmt.Sprintf("duh version %s\nBuilt: %s\nCommit: %s\nGo version: %s\nOS/Arch: %s/%s",
		Version,
		Built,
		Commit,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)
}

// GetVersion returns just the version string
func GetVersion() string {
	return Version
}
