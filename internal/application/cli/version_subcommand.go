package cli

import (
	"duh/internal/domain/utils/version"
	"fmt"

	"github.com/spf13/cobra"
)

// BuildVersionCommand creates the version subcommand
func BuildVersionCommand() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  "Display the version of duh along with build information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.BuildInfo())
		},
	}

	return versionCmd
}
