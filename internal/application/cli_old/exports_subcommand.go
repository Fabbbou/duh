package cli_old

import (
	"duh/internal/domain/service"

	"github.com/spf13/cobra"
)

func BuildExportsSubcommand(cliService service.CliService) *cobra.Command {

	exportsCmd := &cobra.Command{
		Use:     "exports [subcommand]",
		Aliases: []string{"export", "ex"},
		Short:   "Keep exports in your shell for good, duh.",
	}

	setExportCmd := &cobra.Command{
		Use:   "set [export_name] [value]",
		Short: "Set an export for an environment variable",
		Args:  cobra.ExactArgs(2),

		Run: func(cmd *cobra.Command, args []string) {
			exportName := args[0]
			value := args[1]
			err := cliService.UpsertExport(exportName, value)
			if err != nil {
				cmd.PrintErrf("Error setting export: %v\n", err)
				return
			}
			cmd.Printf("Export '%s' set for value '%s'\n", exportName, value)
		},
	}

	unsetExportCmd := &cobra.Command{
		Use:   "unset [export_name]",
		Short: "Remove an export",
		Args:  cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			exportName := args[0]
			err := cliService.RemoveExport(exportName)
			if err != nil {
				cmd.PrintErrf("Error removing export: %v\n", err)
				return
			}
			cmd.Printf("Export '%s' removed\n", exportName)
		},
	}

	listExportCmd := &cobra.Command{
		Use:   "list",
		Short: "List all exports",
		Args:  cobra.NoArgs,

		Run: func(cmd *cobra.Command, args []string) {
			entries, err := cliService.ListExports()
			if err != nil {
				cmd.PrintErrf("%s: %v\n", "Error listing exports", err)
				return
			}
			for key, value := range entries {
				cmd.Printf("%s='%s'\n", key, value)
			}
		},
	}

	exportsCmd.AddCommand(setExportCmd)
	exportsCmd.AddCommand(unsetExportCmd)
	exportsCmd.AddCommand(listExportCmd)

	return exportsCmd
}
