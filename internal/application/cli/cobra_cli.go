package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func BuildRootCli() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "duh",
		Short: "Duh, a simple and effective dotfiles manager",
	}

	aliasRootCmd := &cobra.Command{
		Use:   "alias [subcommand]",
		Short: "Alias management commands",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Hello, %s!\n", args[0])
		},
	}
	rootCmd.AddCommand(aliasRootCmd)

	return rootCmd
}
