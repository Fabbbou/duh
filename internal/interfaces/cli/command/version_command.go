package command

import (
	"duh/internal/interfaces/cli/handler"

	"github.com/spf13/cobra"
)

func BuildVersionCommand(versionHandler *handler.VersionHandler) *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run:   versionHandler.ShowVersion,
	}

	return versionCmd
}
