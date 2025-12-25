package cli

import (
	"duh/internal/domain/service"

	"github.com/spf13/cobra"
)

func BuildAliasSubcommand(cliService service.CliService) *cobra.Command {

	aliasCmd := &cobra.Command{
		Use:   "alias [subcommand]",
		Short: "Keep alias in your shell for good, duh.",
	}

	setAliasCmd := &cobra.Command{
		Use:   "set [alias_name] [value]",
		Short: "Set an alias for a command",
		Args:  cobra.ExactArgs(2),

		Run: func(cmd *cobra.Command, args []string) {
			aliasName := args[0]
			value := args[1]
			err := cliService.UpsertAlias(aliasName, value)
			if err != nil {
				cmd.PrintErrf("Error setting alias: %v\n", err)
				return
			}
			cmd.Printf("Alias '%s' set for command '%s'\n", aliasName, value)
		},
	}

	unsetAliasCmd := &cobra.Command{
		Use:   "unset [alias_name]",
		Short: "Remove an alias",
		Args:  cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			aliasName := args[0]
			err := cliService.RemoveAlias(aliasName)
			if err != nil {
				cmd.PrintErrf("Error removing alias: %v\n", err)
				return
			}
			cmd.Printf("Alias '%s' removed\n", aliasName)
		},
	}

	listAliasCmd := &cobra.Command{
		Use:   "list",
		Short: "List all aliases",
		Args:  cobra.NoArgs,

		Run: func(cmd *cobra.Command, args []string) {
			entries, err := cliService.ListAliases()
			if err != nil {
				cmd.PrintErrf("%s: %v\n", "Error listing aliases", err)
				return
			}
			for key, value := range entries {
				cmd.Printf("%s='%s'\n", key, value)
			}
		},
	}

	aliasCmd.AddCommand(setAliasCmd)
	aliasCmd.AddCommand(unsetAliasCmd)
	aliasCmd.AddCommand(listAliasCmd)

	return aliasCmd
}
