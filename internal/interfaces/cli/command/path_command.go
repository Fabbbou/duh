package command

import (
	"duh/internal/interfaces/cli/handler"

	"github.com/spf13/cobra"
)

func BuildPathCommand(pathHandler *handler.PathHandler) *cobra.Command {
	pathCmd := &cobra.Command{
		Use:     "path [subcommand]",
		Aliases: []string{"paths", "pa"},
		Short:   "Manage and view repository paths",
		Args:    cobra.RangeArgs(0, 1),
		Run:     pathHandler.ShowPath,
	}

	allPathsList := &cobra.Command{
		Use:   "list",
		Short: "List all repository base paths",
		Args:  cobra.NoArgs,
		Run:   pathHandler.ListAllPaths,
	}

	pathCmd.AddCommand(allPathsList)

	return pathCmd
}
