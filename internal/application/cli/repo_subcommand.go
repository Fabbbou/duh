package cli

import (
	"duh/internal/domain/service"

	"github.com/spf13/cobra"
)

func BuildRepoSubcommand(cliService service.CliService) *cobra.Command {
	repoCmd := &cobra.Command{
		Use:     "repository [subcommand]",
		Aliases: []string{"repo", "repos"},
		Short:   "Manage repositories for aliases and exports",
	}

	listRepoCmd := &cobra.Command{
		Use:   "list",
		Short: "List all repositories",
		Args:  cobra.NoArgs,

		Run: func(cmd *cobra.Command, args []string) {
			repos, err := cliService.ListRepositories()
			if err != nil {
				cmd.PrintErrf("Error listing repositories: %v\n", err)
				return
			}

			if len(repos["enabled"]) > 0 {
				cmd.Println("Enabled repositories:")
				for _, name := range repos["enabled"] {
					cmd.Printf("  ✓ %s\n", name)
				}
			}

			if len(repos["disabled"]) > 0 {
				cmd.Println("Disabled repositories:")
				for _, name := range repos["disabled"] {
					cmd.Printf("  ✗ %s\n", name)
				}
			}

			if len(repos["enabled"]) == 0 && len(repos["disabled"]) == 0 {
				cmd.Println("No repositories found")
			}
		},
	}

	enableRepoCmd := &cobra.Command{
		Use:   "enable [repo_name]",
		Short: "Enable a repository",
		Args:  cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			repoName := args[0]
			err := cliService.EnableRepository(repoName)
			if err != nil {
				cmd.PrintErrf("Error enabling repository: %v\n", err)
				return
			}
			cmd.Printf("Repository '%s' enabled\n", repoName)
		},
	}

	disableRepoCmd := &cobra.Command{
		Use:   "disable [repo_name]",
		Short: "Disable a repository",
		Args:  cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			repoName := args[0]
			err := cliService.DisableRepository(repoName)
			if err != nil {
				cmd.PrintErrf("Error disabling repository: %v\n", err)
				return
			}
			cmd.Printf("Repository '%s' disabled\n", repoName)
		},
	}

	deleteRepoCmd := &cobra.Command{
		Use:   "delete [repo_name]",
		Short: "Delete a repository",
		Args:  cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			repoName := args[0]
			err := cliService.DeleteRepository(repoName)
			if err != nil {
				cmd.PrintErrf("Error deleting repository: %v\n", err)
				return
			}
			cmd.Printf("Repository '%s' deleted\n", repoName)
		},
	}

	setDefaultRepoCmd := &cobra.Command{
		Use:   "default [subcommand]",
		Short: "Manage default repository",
		Run: func(cmd *cobra.Command, args []string) {
			// Show current default and help when no subcommand provided
			currentDefault, err := cliService.GetCurrentDefaultRepository()
			if err != nil {
				cmd.PrintErrf("Error getting current default repository: %v\n", err)
				return
			}
			cmd.Printf("Current default repository: %s\n\n", currentDefault)
			cmd.Println("Available commands:")
			cmd.Println("  duh repo default set <name>  Set repository as default")
		},
	}

	setDefaultRepoSetCmd := &cobra.Command{
		Use:   "set [repo_name]",
		Short: "Set a repository as default",
		Args:  cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			repoName := args[0]
			err := cliService.SetDefaultRepository(repoName)
			if err != nil {
				cmd.PrintErrf("Error setting default repository: %v\n", err)
				return
			}
			cmd.Printf("Repository '%s' set as default\n", repoName)
		},
	}

	setDefaultRepoCmd.AddCommand(setDefaultRepoSetCmd)

	renameRepoCmd := &cobra.Command{
		Use:   "rename [old_name] [new_name]",
		Short: "Rename a repository",
		Args:  cobra.ExactArgs(2),

		Run: func(cmd *cobra.Command, args []string) {
			oldName := args[0]
			newName := args[1]
			err := cliService.RenameRepository(oldName, newName)
			if err != nil {
				cmd.PrintErrf("Error renaming repository: %v\n", err)
				return
			}
			cmd.Printf("Repository '%s' renamed to '%s'\n", oldName, newName)
		},
	}

	addRepoCmd := &cobra.Command{
		Use:   "add [url] [name (optional)]",
		Short: "Add a new repository",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			var name *string
			if len(args) == 2 {
				name = &args[1]
			}
			err := cliService.AddRepository(url, name)
			if err != nil {
				cmd.PrintErrf("Error adding repository: %v\n", err)
				return
			}
			cmd.Printf("Repository '%s' added and enabled\n", url)
		},
	}

	createRepoCmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new empty repository",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			repoName := args[0]
			err := cliService.CreateRepository(repoName)
			if err != nil {
				cmd.PrintErrf("Error creating repository: %v\n", err)
				return
			}
			cmd.Printf("Repository '%s' created and enabled\n", repoName)
		},
	}

	repoCmd.AddCommand(listRepoCmd)
	repoCmd.AddCommand(enableRepoCmd)
	repoCmd.AddCommand(disableRepoCmd)
	repoCmd.AddCommand(deleteRepoCmd)
	repoCmd.AddCommand(setDefaultRepoCmd)
	repoCmd.AddCommand(renameRepoCmd)
	repoCmd.AddCommand(addRepoCmd)
	repoCmd.AddCommand(createRepoCmd)

	return repoCmd
}
