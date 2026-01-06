package command

import (
	"duh/internal/interfaces/cli/handler"

	"github.com/spf13/cobra"
)

func BuildSelfCommand(selfHandler *handler.SelfHandler) *cobra.Command {

	selfCmd := &cobra.Command{
		Use:   "self [subcommand]",
		Short: "Commands to manage duh itself.",
	}

	configPathCmd := &cobra.Command{
		Use:   "config-path",
		Short: "Show the path to the configuration directory",
		Args:  cobra.NoArgs,
		Run:   selfHandler.ShowConfigPath,
	}

	repositoriesPathCmd := &cobra.Command{
		Use:   "repositories-path",
		Short: "Show the path to the repositories directory",
		Args:  cobra.NoArgs,
		Run:   selfHandler.ShowRepositoriesPath,
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show the version of duh",
		Args:  cobra.NoArgs,
		Run:   selfHandler.GetVersion,
	}

	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update duh to the latest version",
		Args:  cobra.NoArgs,
		Run:   selfHandler.Update,
	}

	selfCmd.AddCommand(configPathCmd)
	selfCmd.AddCommand(repositoriesPathCmd)
	selfCmd.AddCommand(versionCmd)
	selfCmd.AddCommand(updateCmd)

	return selfCmd
}
