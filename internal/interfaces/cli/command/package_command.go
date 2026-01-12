package command

import (
	"duh/internal/interfaces/cli/handler"

	"github.com/spf13/cobra"
)

func BuildPackageCommand(packageHandler *handler.PackageHandler) *cobra.Command {
	packageCmd := &cobra.Command{
		Use:     "package [subcommand]",
		Aliases: []string{"pkg", "packages", "repo", "repos", "repository"},
		Short:   "Manage packages for aliases and exports",
	}

	listPackageCmd := &cobra.Command{
		Use:   "list",
		Short: "List all packages",
		Args:  cobra.NoArgs,
		Run:   packageHandler.ListPackages,
	}

	enablePackageCmd := &cobra.Command{
		Use:   "enable [package_name]",
		Short: "Enable a package",
		Args:  cobra.ExactArgs(1),
		Run:   packageHandler.EnablePackage,
	}

	disablePackageCmd := &cobra.Command{
		Use:   "disable [package_name]",
		Short: "Disable a package",
		Args:  cobra.ExactArgs(1),
		Run:   packageHandler.DisablePackage,
	}

	deletePackageCmd := &cobra.Command{
		Use:   "delete [package_name]",
		Short: "Delete a package",
		Args:  cobra.ExactArgs(1),
		Run:   packageHandler.DeletePackage,
	}

	// Default package management with subcommands
	setDefaultPackageCmd := &cobra.Command{
		Use:   "default [subcommand]",
		Short: "Manage default package",
		Run:   packageHandler.ShowDefaultPackage,
	}

	setDefaultPackageSetCmd := &cobra.Command{
		Use:   "set [package_name]",
		Short: "Set a package as default",
		Args:  cobra.ExactArgs(1),
		Run:   packageHandler.SetDefaultPackage,
	}

	setDefaultPackageCmd.AddCommand(setDefaultPackageSetCmd)

	getDefaultPackageCmd := &cobra.Command{
		Use:   "current",
		Short: "Show current default package",
		Args:  cobra.NoArgs,
		Run:   packageHandler.GetDefaultPackage,
	}

	renamePackageCmd := &cobra.Command{
		Use:   "rename [old_name] [new_name]",
		Short: "Rename a package",
		Args:  cobra.ExactArgs(2),
		Run:   packageHandler.RenamePackage,
	}

	addPackageCmd := &cobra.Command{
		Use:   "add [url] [name (optional)]",
		Short: "Add a new package",
		Args:  cobra.RangeArgs(1, 2),
		Run:   packageHandler.AddPackage,
	}

	createPackageCmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new empty package",
		Args:  cobra.ExactArgs(1),
		Run:   packageHandler.CreatePackage,
	}

	updatePackageCmd := &cobra.Command{
		Use:   "update",
		Short: "Update packages from their remote sources",
		Long: `Update all enabled packages that have git remotes.
By default, updates are safe and won't proceed if local changes exist.

Strategies:
  --commit  Commit local changes before updating (safer)
  --force   Discard local changes and force update (destructive)

If neither flag is provided, the update will fail if local changes exist.`,
		Args: cobra.NoArgs,
		Run:  packageHandler.UpdatePackages,
	}

	editPackageCmd := &cobra.Command{
		Use:   "edit [package_name]",
		Short: "Edit a package's configuration file",
		Long: `Open the configuration file (db.toml) of the specified package in the system's default text editor.
This allows you to modify package settings such as aliases and exports.

This command tries to use the default editor set in your system.
You can override the default editor by setting the $EDITOR environment variable.
For example:
duh export EDITOR nano
`,
		Args: cobra.ExactArgs(1),
		Run:  packageHandler.EditPackage,
	}

	pushPackageCmd := &cobra.Command{
		Use:   "push [package_name]",
		Short: "Push local changes to remote package",
		Long: `Push local changes to the remote package. If there are uncommitted changes,
they will be automatically committed before pushing.

This command requires:
- The package must have a git remote configured
- You must have push permissions to the remote package

Example:
  duh package push my-package
`,
		Args: cobra.ExactArgs(1),
		Run:  packageHandler.PushPackage,
	}

	editGitconfigCmd := &cobra.Command{
		Use:     "edit-gitconfig [package_name]",
		Aliases: []string{"git"},
		Short:   "Edit a package's gitconfig file",
		Long: `Open the gitconfig file of the specified package in the system's default text editor.
This allows you to modify git-specific settings for the package.

This command tries to use the default editor set in your system.
You can override the default editor by setting the $EDITOR environment variable.
For example:
duh export EDITOR nano
`,
		Args: cobra.ExactArgs(1),
		Run:  packageHandler.EditGitconfig,
	}

	// Add flags to update command
	updatePackageCmd.Flags().Bool("force", false, "Force update by discarding local changes (destructive)")
	updatePackageCmd.Flags().Bool("commit", false, "Commit local changes before updating (safer)")

	packageCmd.AddCommand(listPackageCmd)
	packageCmd.AddCommand(enablePackageCmd)
	packageCmd.AddCommand(disablePackageCmd)
	packageCmd.AddCommand(deletePackageCmd)
	packageCmd.AddCommand(setDefaultPackageCmd)
	packageCmd.AddCommand(getDefaultPackageCmd)
	packageCmd.AddCommand(renamePackageCmd)
	packageCmd.AddCommand(addPackageCmd)
	packageCmd.AddCommand(createPackageCmd)
	packageCmd.AddCommand(updatePackageCmd)
	packageCmd.AddCommand(editPackageCmd)
	packageCmd.AddCommand(editGitconfigCmd)
	packageCmd.AddCommand(pushPackageCmd)

	return packageCmd
}
