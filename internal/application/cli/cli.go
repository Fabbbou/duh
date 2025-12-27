package cli

import (
	"duh/internal/domain/service"
	"duh/internal/version"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func BuildRootCli(cliService service.CliService) *cobra.Command {
	var showVersion bool

	rootCmd := &cobra.Command{
		Use:   "duh",
		Short: "Duh, a simple and effective dotfiles manager",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			CheckDuhFileDBCreated(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			if showVersion {
				fmt.Println("duh", version.GetVersion())
				os.Exit(0)
			} else {
				cmd.Help()
			}
		},
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	// Add global --version flag
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Show version information")

	// Add version subcommand
	versionCmd := BuildVersionCommand()
	rootCmd.AddCommand(versionCmd)

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
