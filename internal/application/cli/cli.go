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

	// exportsCmd := &cobra.Command{
	// 	Use:   "exports [subcommand]",
	// 	Short: "in your shell for good, duh",
	// 	Args:  cobra.ExactArgs(1),

	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		fmt.Printf("Hello, %s!\n", args[0])
	// 	},
	// }

	// selfCmd := &cobra.Command{
	// 	Use:   "self [subcommand]",
	// 	Short: "Duh, itself.",
	// 	Args:  cobra.ExactArgs(1),

	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		fmt.Printf("Hello, %s!\n", args[0])
	// 	},
	// }

	rootCmd.AddCommand(aliasCmd)
	// rootCmd.AddCommand(exportsCmd)
	// rootCmd.AddCommand(selfCmd)

	return rootCmd
}
