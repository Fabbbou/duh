package command

import (
	"duh/internal/interfaces/cli/handler"

	"github.com/spf13/cobra"
)

func BuildRootCli(
	initFileDBHandler *handler.InitFileDBHandler,
	aliasHandler *handler.AliasHandler,
	exportsHandler *handler.ExportsHandler,
	functionsHandler *handler.FunctionsHandler,
	injectHandler *handler.InjectHandler,
	packageHandler *handler.PackageHandler,
	selfHandler *handler.SelfHandler,
) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "duh",
		Short: "Duh, a simple and effective dotfiles manager",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initFileDBHandler.HandleInitFileDB(cmd)
		},
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	// Add all subcommands using new handlers
	rootCmd.AddCommand(BuildInjectCommand(injectHandler))
	rootCmd.AddCommand(BuildAliasCommand(aliasHandler))
	rootCmd.AddCommand(BuildExportsCommand(exportsHandler))
	rootCmd.AddCommand(BuildPackageCommand(packageHandler))
	rootCmd.AddCommand(BuildFunctionsCommand(functionsHandler))
	rootCmd.AddCommand(BuildSelfCommand(selfHandler))

	return rootCmd
}
