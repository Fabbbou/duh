package cli

import (
	"duh/internal/domain/service"

	"github.com/spf13/cobra"
)

func BuildRootCli(cliService service.CliService) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "duh",
		Short: "Duh, a simple and effective dotfiles manager",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			CheckDuhFileDBCreated(cmd)
		},
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	injectCmd := BuildInjectSubcommand(cliService)
	rootCmd.AddCommand(injectCmd)

	aliasCmd := BuildAliasSubcommand(cliService)
	rootCmd.AddCommand(aliasCmd)

	exportsCmd := BuildExportsSubcommand(cliService)
	rootCmd.AddCommand(exportsCmd)

	repoCmd := BuildRepoSubcommand(cliService)
	rootCmd.AddCommand(repoCmd)

	pathCmd := BuildPathSubcommand(cliService)
	rootCmd.AddCommand(pathCmd)

	return rootCmd
}
