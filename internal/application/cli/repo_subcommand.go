package cli

import (
	"duh/internal/domain/entity"
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

	updateRepoCmd := &cobra.Command{
		Use:   "update",
		Short: "Update repositories from their remote sources",
		Long: `Update all enabled repositories that have git remotes.
By default, updates are safe and won't proceed if local changes exist.

Strategies:
  --commit  Commit local changes before updating (safer)
  --force   Discard local changes and force update (destructive)

If neither flag is provided, the update will fail if local changes exist.`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// Determine strategy based on flags
			forceFlag, _ := cmd.Flags().GetBool("force")
			commitFlag, _ := cmd.Flags().GetBool("commit")

			strategy := entity.UpdateSafe // default strategy
			if forceFlag && commitFlag {
				cmd.PrintErrf("Cannot use both --force and --commit flags together\n")
				return
			} else if forceFlag {
				strategy = entity.UpdateForce
			} else if commitFlag {
				strategy = entity.UpdateKeep
			}

			results, err := cliService.UpdateRepos(strategy)
			if err != nil {
				cmd.PrintErrf("Error updating repositories: %v\n", err)
				return
			}

			// Report results
			if len(results.LocalChangesDetected) > 0 {
				cmd.Println("⚠️  Repositories with local changes (not updated):")
				for _, repo := range results.LocalChangesDetected {
					cmd.Printf("  • %s\n", repo)
				}
				cmd.Println("\nUse --commit to auto-commit changes or --force to discard them.")
			}

			if len(results.OtherErrors) > 0 {
				cmd.Println("❌ Repositories with errors:")
				for _, err := range results.OtherErrors {
					cmd.Printf("  • %v\n", err)
				}
			}

			totalRepos := len(results.LocalChangesDetected) + len(results.OtherErrors)
			if totalRepos == 0 {
				cmd.Println("✅ All repositories updated successfully")
			} else {
				cmd.Printf("\n%d repositories had issues during update\n", totalRepos)
			}
		},
	}

	editRepoCmd := &cobra.Command{
		Use:   "edit [repo_name]",
		Short: "Edit a repository's configuration file",
		Long: `Open the configuration file (db.toml) of the specified repository in the system's default text editor.
This allows you to modify repository settings such as aliases and exports.

This command tries to use the default editor set in your system.
You can override the default editor by setting the $EDITOR environment variable.
For example:
duh export EDITOR nano
`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			repoName := args[0]
			err := cliService.EditRepo(repoName)
			if err != nil {
				cmd.PrintErrf("Error editing repository: %v\n", err)
				return
			}
		},
	}

	pushRepoCmd := &cobra.Command{
		Use:   "push [repo_name]",
		Short: "Push local changes to remote repository",
		Long: `Push local changes to the remote repository. If there are uncommitted changes,
they will be automatically committed before pushing.

This command requires:
- The repository must have a git remote configured
- You must have push permissions to the remote repository

Example:
  duh repo push my-repo
`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			repoName := args[0]
			err := cliService.PushRepository(repoName)
			if err != nil {
				cmd.PrintErrf("Error pushing repository: %v\n", err)
				return
			}
			cmd.Printf("Repository '%s' pushed successfully\n", repoName)
		},
	}

	// Add flags to update command
	updateRepoCmd.Flags().Bool("force", false, "Force update by discarding local changes (destructive)")
	updateRepoCmd.Flags().Bool("commit", false, "Commit local changes before updating (safer)")

	repoCmd.AddCommand(listRepoCmd)
	repoCmd.AddCommand(enableRepoCmd)
	repoCmd.AddCommand(disableRepoCmd)
	repoCmd.AddCommand(deleteRepoCmd)
	repoCmd.AddCommand(setDefaultRepoCmd)
	repoCmd.AddCommand(renameRepoCmd)
	repoCmd.AddCommand(addRepoCmd)
	repoCmd.AddCommand(createRepoCmd)
	repoCmd.AddCommand(updateRepoCmd)
	repoCmd.AddCommand(editRepoCmd)
	repoCmd.AddCommand(pushRepoCmd)

	return repoCmd
}
