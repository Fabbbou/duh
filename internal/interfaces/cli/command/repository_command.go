package command

import (
	"duh/internal/interfaces/cli/handler"

	"github.com/spf13/cobra"
)

func BuildRepositoryCommand(repositoryHandler *handler.RepositoryHandler) *cobra.Command {
	repoCmd := &cobra.Command{
		Use:     "repository [subcommand]",
		Aliases: []string{"repo", "repos"},
		Short:   "Manage repositories for aliases and exports",
	}

	listRepoCmd := &cobra.Command{
		Use:   "list",
		Short: "List all repositories",
		Args:  cobra.NoArgs,
		Run:   repositoryHandler.ListRepositories,
	}

	enableRepoCmd := &cobra.Command{
		Use:   "enable [repo_name]",
		Short: "Enable a repository",
		Args:  cobra.ExactArgs(1),
		Run:   repositoryHandler.EnableRepository,
	}

	disableRepoCmd := &cobra.Command{
		Use:   "disable [repo_name]",
		Short: "Disable a repository",
		Args:  cobra.ExactArgs(1),
		Run:   repositoryHandler.DisableRepository,
	}

	deleteRepoCmd := &cobra.Command{
		Use:   "delete [repo_name]",
		Short: "Delete a repository",
		Args:  cobra.ExactArgs(1),
		Run:   repositoryHandler.DeleteRepository,
	}

	// Default repository management with subcommands
	setDefaultRepoCmd := &cobra.Command{
		Use:   "default [subcommand]",
		Short: "Manage default repository",
		Run:   repositoryHandler.ShowDefaultRepository,
	}

	setDefaultRepoSetCmd := &cobra.Command{
		Use:   "set [repo_name]",
		Short: "Set a repository as default",
		Args:  cobra.ExactArgs(1),
		Run:   repositoryHandler.SetDefaultRepository,
	}

	setDefaultRepoCmd.AddCommand(setDefaultRepoSetCmd)

	getDefaultRepoCmd := &cobra.Command{
		Use:   "current",
		Short: "Show current default repository",
		Args:  cobra.NoArgs,
		Run:   repositoryHandler.GetDefaultRepository,
	}

	renameRepoCmd := &cobra.Command{
		Use:   "rename [old_name] [new_name]",
		Short: "Rename a repository",
		Args:  cobra.ExactArgs(2),
		Run:   repositoryHandler.RenameRepository,
	}

	addRepoCmd := &cobra.Command{
		Use:   "add [url] [name (optional)]",
		Short: "Add a new repository",
		Args:  cobra.RangeArgs(1, 2),
		Run:   repositoryHandler.AddRepository,
	}

	createRepoCmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new empty repository",
		Args:  cobra.ExactArgs(1),
		Run:   repositoryHandler.CreateRepository,
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
		Run:  repositoryHandler.UpdateRepositories,
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
		Run:  repositoryHandler.EditRepository,
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
		Run:  repositoryHandler.PushRepository,
	}

	editGitconfigCmd := &cobra.Command{
		Use:     "edit-gitconfig [repo_name]",
		Aliases: []string{"git"},
		Short:   "Edit a repository's gitconfig file",
		Long: `Open the gitconfig file of the specified repository in the system's default text editor.
This allows you to modify git-specific settings for the repository.

This command tries to use the default editor set in your system.
You can override the default editor by setting the $EDITOR environment variable.
For example:
duh export EDITOR nano
`,
		Args: cobra.ExactArgs(1),
		Run:  repositoryHandler.EditGitconfig,
	}

	// Add flags to update command
	updateRepoCmd.Flags().Bool("force", false, "Force update by discarding local changes (destructive)")
	updateRepoCmd.Flags().Bool("commit", false, "Commit local changes before updating (safer)")

	repoCmd.AddCommand(listRepoCmd)
	repoCmd.AddCommand(enableRepoCmd)
	repoCmd.AddCommand(disableRepoCmd)
	repoCmd.AddCommand(deleteRepoCmd)
	repoCmd.AddCommand(setDefaultRepoCmd)
	repoCmd.AddCommand(getDefaultRepoCmd)
	repoCmd.AddCommand(renameRepoCmd)
	repoCmd.AddCommand(addRepoCmd)
	repoCmd.AddCommand(createRepoCmd)
	repoCmd.AddCommand(updateRepoCmd)
	repoCmd.AddCommand(editRepoCmd)
	repoCmd.AddCommand(editGitconfigCmd)
	repoCmd.AddCommand(pushRepoCmd)

	return repoCmd
}
