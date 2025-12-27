package gitt

import "github.com/go-git/go-git/v5"

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
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}
	err = repo.Push(&git.PushOptions{})
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
