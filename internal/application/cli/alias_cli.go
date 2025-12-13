package cli

import (
	"duh/internal/domain/service"

	"github.com/spf13/cobra"
)

func BuildAliasCli(cliService service.CliService) *cobra.Command {

	aliasCmd := &cobra.Command{
		Use:   "alias [subcommand]",
		Short: "Keep alias in your shell for good, duh.",
		Args:  cobra.ExactArgs(1),

		// Run: func(cmd *cobra.Command, args []string) {
		// 	fmt.Printf("Hello, %s!\n", args[0])
		// },
	}

	setAliasCmd := &cobra.Command{
		Use:   "set [alias_name] [command]",
		Short: "Set an alias for a command",
		Args:  cobra.ExactArgs(2),

		Run: func(cmd *cobra.Command, args []string) {
			aliasName := args[0]
			command := args[1]
			err := cliService.SetAlias(aliasName, command)
			if err != nil {
				cmd.PrintErrf("Error setting alias: %v\n", err)
				return
			}
			cmd.Printf("Alias '%s' set for command '%s'\n", aliasName, command)
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
			repo, err := cliService.ListAliases()
			if err != nil {
				cmd.PrintErrf("Error listing aliases: %v\n", err)
				return
			}
			for key, value := range repo.Aliases {
				cmd.Printf("%s='%s'\n", key, value)
			}
		},
	}

	aliasCmd.AddCommand(setAliasCmd)
	aliasCmd.AddCommand(unsetAliasCmd)
	aliasCmd.AddCommand(listAliasCmd)

	return aliasCmd
}
