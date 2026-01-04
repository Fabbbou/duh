package command

import (
	"duh/internal/interfaces/cli/handler"

	"github.com/spf13/cobra"
)

func BuildAliasCommand(aliasHandler *handler.AliasHandler) *cobra.Command {

	aliasCmd := &cobra.Command{
		Use:   "alias [subcommand]",
		Short: "Keep alias in your shell for good, duh.",
	}

	setAliasCmd := &cobra.Command{
		Use:   "set [alias_name] [value]",
		Short: "Set an alias for a command",
		Args:  cobra.ExactArgs(2),
		Run:   aliasHandler.SetAlias,
	}

	unsetAliasCmd := &cobra.Command{
		Use:   "unset [alias_name]",
		Short: "Remove an alias",
		Args:  cobra.ExactArgs(1),
		Run:   aliasHandler.UnsetAlias,
	}

	listAliasCmd := &cobra.Command{
		Use:   "list",
		Short: "List all aliases",
		Args:  cobra.NoArgs,
		Run:   aliasHandler.ListAliases,
	}

	aliasCmd.AddCommand(setAliasCmd)
	aliasCmd.AddCommand(unsetAliasCmd)
	aliasCmd.AddCommand(listAliasCmd)
	return aliasCmd
}
