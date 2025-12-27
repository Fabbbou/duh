package cli

import (
	"duh/internal/version"
	"fmt"

	"github.com/spf13/cobra"
)

// BuildVersionCommand creates the version subcommand
func BuildVersionCommand() *cobra.Command {
	var detailed bool

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  "Display the version of duh along with build information",
		Run: func(cmd *cobra.Command, args []string) {
			if detailed {
				fmt.Println(version.BuildInfo())
			} else {
				fmt.Println("duh", version.GetVersion())
			}
		},
	}

	versionCmd.Flags().BoolVarP(&detailed, "detailed", "d", false, "Show detailed build information")

	return versionCmd
}
