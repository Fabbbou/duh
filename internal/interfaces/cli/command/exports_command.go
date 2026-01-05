package command

import (
	"duh/internal/interfaces/cli/handler"

	"github.com/spf13/cobra"
)

func BuildExportsCommand(exportsHandler *handler.ExportsHandler) *cobra.Command {
	exportsCmd := &cobra.Command{
		Use:     "exports [subcommand]",
		Aliases: []string{"export", "ex"},
		Short:   "Keep exports in your shell for good, duh.",
	}

	setExportCmd := &cobra.Command{
		Use:   "set [export_name] [value]",
		Short: "Set an export for an environment variable",
		Args:  cobra.ExactArgs(2),
		Run:   exportsHandler.SetExport,
	}

	unsetExportCmd := &cobra.Command{
		Use:   "unset [export_name]",
		Short: "Remove an export",
		Args:  cobra.ExactArgs(1),
		Run:   exportsHandler.UnsetExport,
	}

	listExportCmd := &cobra.Command{
		Use:   "list",
		Short: "List all exports",
		Args:  cobra.NoArgs,
		Run:   exportsHandler.ListExports,
	}

	exportsCmd.AddCommand(setExportCmd)
	exportsCmd.AddCommand(unsetExportCmd)
	exportsCmd.AddCommand(listExportCmd)

	return exportsCmd
}
