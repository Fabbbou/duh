package handler

import (
	"duh/internal/application/usecase"
	"duh/internal/domain/entity"

	"github.com/spf13/cobra"
)

type RepositoryHandler struct {
	repositoryUsecase *usecase.RepositoryUsecase
}

func NewRepositoryHandler(repositoryUsecase *usecase.RepositoryUsecase) *RepositoryHandler {
	return &RepositoryHandler{
		repositoryUsecase: repositoryUsecase,
	}
}

func (r *RepositoryHandler) ListRepositories(cmd *cobra.Command, args []string) {
	repos, err := r.repositoryUsecase.ListRepositories()
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
}

func (r *RepositoryHandler) EnableRepository(cmd *cobra.Command, args []string) {
	repoName := args[0]
	err := r.repositoryUsecase.EnableRepository(repoName)
	if err != nil {
		cmd.PrintErrf("Error enabling repository: %v\n", err)
		return
	}
	cmd.Printf("Repository '%s' enabled\n", repoName)
}

func (r *RepositoryHandler) DisableRepository(cmd *cobra.Command, args []string) {
	repoName := args[0]
	err := r.repositoryUsecase.DisableRepository(repoName)
	if err != nil {
		cmd.PrintErrf("Error disabling repository: %v\n", err)
		return
	}
	cmd.Printf("Repository '%s' disabled\n", repoName)
}

func (r *RepositoryHandler) DeleteRepository(cmd *cobra.Command, args []string) {
	repoName := args[0]
	err := r.repositoryUsecase.DeleteRepository(repoName)
	if err != nil {
		cmd.PrintErrf("Error deleting repository: %v\n", err)
		return
	}
	cmd.Printf("Repository '%s' deleted\n", repoName)
}

func (r *RepositoryHandler) SetDefaultRepository(cmd *cobra.Command, args []string) {
	repoName := args[0]
	err := r.repositoryUsecase.SetDefaultRepository(repoName)
	if err != nil {
		cmd.PrintErrf("Error setting default repository: %v\n", err)
		return
	}
	cmd.Printf("Default repository set to '%s'\n", repoName)
}

func (r *RepositoryHandler) GetDefaultRepository(cmd *cobra.Command, args []string) {
	repoName, err := r.repositoryUsecase.GetDefaultRepository()
	if err != nil {
		cmd.PrintErrf("Error getting default repository: %v\n", err)
		return
	}
	cmd.Printf("Default repository: %s\n", repoName)
}

func (r *RepositoryHandler) ShowDefaultRepository(cmd *cobra.Command, args []string) {
	// Show current default and help when no subcommand provided
	currentDefault, err := r.repositoryUsecase.GetDefaultRepository()
	if err != nil {
		cmd.PrintErrf("Error getting current default repository: %v\n", err)
		return
	}
	cmd.Printf("Current default repository: %s\n\n", currentDefault)
	cmd.Println("Available commands:")
	cmd.Println("  duh repo default set <name>  Set repository as default")
}

func (r *RepositoryHandler) RenameRepository(cmd *cobra.Command, args []string) {
	oldName := args[0]
	newName := args[1]
	err := r.repositoryUsecase.RenameRepository(oldName, newName)
	if err != nil {
		cmd.PrintErrf("Error renaming repository: %v\n", err)
		return
	}
	cmd.Printf("Repository '%s' renamed to '%s'\n", oldName, newName)
}

func (r *RepositoryHandler) AddRepository(cmd *cobra.Command, args []string) {
	url := args[0]
	var name *string
	if len(args) == 2 {
		name = &args[1]
	}
	err := r.repositoryUsecase.AddRepository(url, name)
	if err != nil {
		cmd.PrintErrf("Error adding repository: %v\n", err)
		return
	}
	cmd.Printf("Repository '%s' added and enabled\n", url)
}

func (r *RepositoryHandler) CreateRepository(cmd *cobra.Command, args []string) {
	repoName := args[0]
	err := r.repositoryUsecase.CreateRepository(repoName)
	if err != nil {
		cmd.PrintErrf("Error creating repository: %v\n", err)
		return
	}
	cmd.Printf("Repository '%s' created and enabled\n", repoName)
}

func (r *RepositoryHandler) UpdateRepositories(cmd *cobra.Command, args []string) {
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

	results, err := r.repositoryUsecase.UpdateRepositories(strategy)
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
}

func (r *RepositoryHandler) EditRepository(cmd *cobra.Command, args []string) {
	repoName := args[0]
	err := r.repositoryUsecase.EditRepository(repoName)
	if err != nil {
		cmd.PrintErrf("Error editing repository: %v\n", err)
		return
	}
}

func (r *RepositoryHandler) PushRepository(cmd *cobra.Command, args []string) {
	repoName := args[0]
	err := r.repositoryUsecase.PushRepository(repoName)
	if err != nil {
		cmd.PrintErrf("Error pushing repository: %v\n", err)
		return
	}
	cmd.Printf("Repository '%s' pushed successfully\n", repoName)
}

func (r *RepositoryHandler) EditGitconfig(cmd *cobra.Command, args []string) {
	repoName := args[0]
	err := r.repositoryUsecase.EditGitconfig(repoName)
	if err != nil {
		cmd.PrintErrf("Error editing gitconfig: %v\n", err)
		return
	}
}
