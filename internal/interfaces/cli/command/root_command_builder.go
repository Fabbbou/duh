package command

import (
	"duh/internal/interfaces/cli/handler"

	"github.com/spf13/cobra"
)

func BuildRootCommand(aliasHandler *handler.AliasHandler) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "duh",
		Short: "duh is a CLI tool to manage your shell aliases easily.",
	}

	rootCmd.AddCommand(BuildAliasCommand(aliasHandler))

	return rootCmd
}
