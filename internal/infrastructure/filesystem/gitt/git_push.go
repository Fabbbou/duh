package gitt

import (
	"fmt"
	"os/exec"

	"github.com/go-git/go-git/v5"
)

func CommitAndPushChanges(repoPath string) error {
	checkWorkingTreeClean, err := checkWorkingTreeClean(repoPath)
	if err != nil {
		return err
	}

	if !checkWorkingTreeClean {
		err := addAndCommitAllChanges(repoPath, "Duh, auto-commit before push")
		if err != nil {
			return err
		}
	}

	// Use git CLI for pushing - it handles authentication automatically
	err = pushUsingGitCLI(repoPath)
	if err != nil {
		return err
	}
	return nil
}

func addAndCommitAllChanges(repoPath string, message string) error {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}
	w, err := repo.Worktree()
	if err != nil {
		return err
	}
	_, err = w.Add(".")
	if err != nil {
		return err
	}
	_, err = w.Commit(message, &git.CommitOptions{})
	if err != nil {
		return err
	}
	return nil
}

// pushUsingGitCLI uses the git CLI to push changes, which handles authentication automatically
func pushUsingGitCLI(repoPath string) error {
	// Change to the repository directory
	cmd := exec.Command("git", "push")
	cmd.Dir = repoPath

	// Run the command and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Include the git command output in the error for better debugging
		return fmt.Errorf("git push failed: %v\nOutput: %s", err, string(output))
	}

	return nil
}
