package command

import (
	"duh/internal/application/cli_old"
	"duh/internal/domain/service"
	"duh/internal/interfaces/cli/handler"

	"github.com/spf13/cobra"
)

//this file will be replaced by root_command_builder.go in the future

func BuildRootCli(
	cliService service.CliService,
	aliasHandler *handler.AliasHandler,
) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "duh",
		Short: "Duh, a simple and effective dotfiles manager",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cli_old.CheckDuhFileDBCreated(cmd)
		},
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	// Add version subcommand
	versionCmd := cli_old.BuildVersionCommand()
	rootCmd.AddCommand(versionCmd)

	injectCmd := cli_old.BuildInjectSubcommand(cliService)
	rootCmd.AddCommand(injectCmd)

	aliasCmd := cli_old.BuildAliasSubcommand(cliService)
	rootCmd.AddCommand(aliasCmd)

	exportsCmd := cli_old.BuildExportsSubcommand(cliService)
	rootCmd.AddCommand(exportsCmd)

	repoCmd := cli_old.BuildRepoSubcommand(cliService)
	rootCmd.AddCommand(repoCmd)

	pathCmd := cli_old.BuildPathSubcommand(cliService)
	rootCmd.AddCommand(pathCmd)

	functionsCmd := cli_old.BuildFunctionsSubcommand(cliService)
	rootCmd.AddCommand(functionsCmd)

	return rootCmd
}
