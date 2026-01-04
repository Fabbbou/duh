package cli_old

import (
	"duh/internal/domain/service"

	"github.com/spf13/cobra"
)

func BuildPathSubcommand(cliService service.CliService) *cobra.Command {
	pathCmd := &cobra.Command{
		Use:     "path [subcommand]",
		Aliases: []string{"paths", "pa"},
		Short:   "Manage and view repository paths",
		Args:    cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			path, err := cliService.GetBasePath()
			if err != nil {
				cmd.PrintErrf("Error retrieving repository paths: %v\n", err)
				return
			}
			stdPrint(path)
		},
	}

	allPathsList := &cobra.Command{
		Use:   "list",
		Short: "List all repository base paths",
		Args:  cobra.NoArgs,

		Run: func(cmd *cobra.Command, args []string) {
			basePath, err := cliService.GetBasePath()
			if err != nil {
				cmd.PrintErrf("Error retrieving Duh path: %v\n", err)
				return
			}
			stdPrint(basePath + "\n")

			paths, err := cliService.ListPath()
			if err != nil {
				cmd.PrintErrf("Error retrieving repository paths: %v\n", err)
				return
			}
			for _, repoPath := range paths {
				stdPrint(repoPath + "\n")
			}
		},
	}

	pathCmd.AddCommand(allPathsList)

	return pathCmd
}
