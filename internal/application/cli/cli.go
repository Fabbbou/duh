package cli

import (
	"duh/internal/domain/service"

	"github.com/spf13/cobra"
)

func BuildRootCli(cliService service.CliService) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "duh",
		Short: "Duh, a simple and effective dotfiles manager",
	}

	aliasCmd := BuildAliasCli(cliService)
	exportsCmd := BuildExportsCli(cliService)

	rootCmd.AddCommand(aliasCmd)
	rootCmd.AddCommand(exportsCmd)

	return rootCmd
}
