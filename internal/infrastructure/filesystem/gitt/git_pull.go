package gitt

import (
	"duh/internal/domain/entity"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Check whether there are updates pending to be pulled from the remote repository
func hasUpdatesPending(repoPath string) (bool, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return false, err
	}

	// Fetch latest changes from remote to update remote references
	err = repo.Fetch(&git.FetchOptions{})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return false, err
	}

	// Get local HEAD commit hash
	localRef, err := repo.Head()
	if err != nil {
		return false, err
	}
	localCommitHash := localRef.Hash()

	// Get remote HEAD commit hash
	remoteRef, err := repo.Reference("refs/remotes/origin/HEAD", true)
	if err != nil {
		// If remote HEAD doesn't exist, try to get the default branch
		// Get the current branch name and check its remote counterpart
		branchName := localRef.Name().Short()
		remoteRefName := "refs/remotes/origin/" + branchName
		remoteRef, err = repo.Reference(plumbing.ReferenceName(remoteRefName), true)
		if err != nil {
			return false, err
		}
	}
	remoteCommitHash := remoteRef.Hash()

	// Compare local and remote commit hashes
	return localCommitHash != remoteCommitHash, nil
}

func updateRepositoryFromRemote(repoPath string) error {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}
	w, err := repo.Worktree()
	if err != nil {
		return err
	}
	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}
	return nil
}

// checkWorkingTreeClean checks if the working tree has any uncommitted changes
func checkWorkingTreeClean(repoPath string) (bool, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return false, err
	}

	w, err := repo.Worktree()
	if err != nil {
		return false, err
	}

	status, err := w.Status()
	if err != nil {
		return false, err
	}

	return status.IsClean(), nil
}

// Pulls changes from remote, handling local changes based on the specified strategy
// Strategies:
// - entity.UpdateSafe: Do not pull if local changes exist, return ErrChangesExist if changes are present
// - entity.UpdateKeep: Commit local changes before pulling
// - entity.UpdateForce: Discard local changes and reset to remote state
func pullWithLocalChanges(repoPath string, strategy string) error {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	// Check if working tree is clean
	isClean, err := checkWorkingTreeClean(repoPath)
	if err != nil {
		return err
	}

	if isClean {
		// No local changes, just do a normal pull
		return updateRepositoryFromRemote(repoPath)
	}

	// Handle local changes based on strategy
	switch strategy {
	case entity.UpdateKeep:
		// Commit local changes first, then pull
		_, err = worktree.Add(".")
		if err != nil {
			return fmt.Errorf("failed to stage changes: %w", err)
		}

		_, err = worktree.Commit("Auto-commit before pull", &git.CommitOptions{
			Author: &object.Signature{
				Name:  "duh",
				Email: "duh@localhost",
				When:  time.Now(),
			},
		})
		if err != nil {
			return fmt.Errorf("failed to commit changes: %w", err)
		}

		// Now try to pull
		err = worktree.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return fmt.Errorf("failed to pull after commit: %w", err)
		}

		return nil

	case entity.UpdateForce:
		// Reset to remote HEAD (WARNING: this will lose local changes)
		err = repo.Fetch(&git.FetchOptions{})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return fmt.Errorf("failed to fetch: %w", err)
		}

		// Get remote HEAD with fallback to current branch's remote counterpart
		remoteRef, err := repo.Reference("refs/remotes/origin/HEAD", true)
		if err != nil {
			// Try getting the current branch's remote counterpart
			head, err := repo.Head()
			if err != nil {
				return fmt.Errorf("failed to get HEAD: %w", err)
			}
			branchName := head.Name().Short()
			remoteRefName := "refs/remotes/origin/" + branchName
			remoteRef, err = repo.Reference(plumbing.ReferenceName(remoteRefName), true)
			if err != nil {
				return fmt.Errorf("failed to get remote reference: %w", err)
			}
		}

		err = worktree.Reset(&git.ResetOptions{
			Commit: remoteRef.Hash(),
			Mode:   git.HardReset,
		})

		if err != nil {
			return fmt.Errorf("failed to reset to remote: %w", err)
		}

		return nil

	case entity.UpdateSafe:
		// Do not pull and return an error indicating local changes exist
		return ErrChangesExist
	default:
		return fmt.Errorf("unknown strategy '%s'. Available strategies: 'keep', 'force', 'safe'", strategy)
	}
}

func hasGitRepoLinked(repoPath string) (bool, error) {
	repo, err := git.PlainOpen(repoPath)
	if err == git.ErrRepositoryNotExists {
		return false, nil
	}
	remotes, err := repo.Remotes()
	if len(remotes) == 0 {
		return false, nil
	}
	return true, nil
}

// Returns a list of all repositories that have a git remote linked
func getAllRepositoryWithRemotes(repoBasePath string) ([]string, error) {
	dirs, err := os.ReadDir(repoBasePath)
	if err != nil {
		return nil, err
	}
	var repos []string
	for _, dir := range dirs {
		if dir.IsDir() {
			repoPath := repoBasePath + string(os.PathSeparator) + dir.Name()
			hasGit, err := hasGitRepoLinked(repoPath)
			if err != nil {
				return nil, err
			}
			if hasGit {
				repos = append(repos, dir.Name())
			}
		}
	}
	return repos, nil
}

// PullAllRepositories pulls updates for all git repositories found in the specified base path
// using the specified strategy for handling local changes.
// Strategies:
// - entity.UpdateSafe: Do not pull if local changes exist, return ErrChangesExist if changes are present
// - entity.UpdateKeep: Commit local changes before pulling
// - entity.UpdateForce: Discard local changes and reset to remote state
func PullAllRepositories(repoBasePath string, strategy string) (entity.PackageUpdateResults, error) {
	repos, err := getAllRepositoryWithRemotes(repoBasePath)
	if err != nil {
		return entity.PackageUpdateResults{}, err
	}

	localChangesDetected := []string{}
	otherErrors := []error{}
	for _, repoName := range repos {
		repoPath := filepath.Join(repoBasePath, repoName)
		err := pullWithLocalChanges(repoPath, strategy)
		if err == ErrChangesExist {
			localChangesDetected = append(localChangesDetected, repoName)
		} else if err != nil {
			otherErrors = append(otherErrors, fmt.Errorf("failed to pull repository '%s': %w", repoName, err))
		}
	}
	return entity.PackageUpdateResults{
		LocalChangesDetected: localChangesDetected,
		OtherErrors:          otherErrors,
	}, nil
}
