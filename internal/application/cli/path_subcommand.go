package cli

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
			cmd.Printf("%s", path)
		},
	}

	allPathsList := &cobra.Command{
		Use:   "list",
		Short: "List all repository base paths",
		Args:  cobra.NoArgs,

		Run: func(cmd *cobra.Command, args []string) {
			basePath, err := cliService.GetBasePath()
			if err != nil {
				cmd.PrintErrf("Error retrieving repository paths: %v\n", err)
				return
			}
			cmd.Printf("%s", basePath)

			paths, err := cliService.ListPath()
			if err != nil {
				cmd.PrintErrf("Error retrieving repository paths: %v\n", err)
				return
			}
			for _, repoPath := range paths {
				cmd.Println(repoPath)
			}
		},
	}

	pathCmd.AddCommand(allPathsList)

	return pathCmd
}
