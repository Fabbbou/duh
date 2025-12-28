package version

import (
	"fmt"
	"runtime"
)

// These variables will be set at build time using ldflags
var (
	// Version is the current version of the application
	Version = "v0.3.3"
)

// BuildInfo returns detailed build information
func BuildInfo() string {
	return fmt.Sprintf("duh version %s\nGo version: %s\nOS/Arch: %s/%s",
		Version,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)
}

// GetVersion returns just the version string
func GetVersion() string {
	return Version
}
