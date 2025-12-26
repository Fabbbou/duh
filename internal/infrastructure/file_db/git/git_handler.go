package gitt

import (
	"os"

	"github.com/go-git/go-git/v5"
)

func CloneGitRepository(url string, outputPath string) error {
	_, err := git.PlainClone(outputPath, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		return err
	}
	return nil
}

func ExtractGitRepoName(url string) string {
	// This is a simplified extraction logic; in real scenarios, consider edge cases.
	base := url
	if len(base) == 0 {
		return ""
	}
	// Remove trailing slash if present
	if base[len(base)-1] == '/' {
		base = base[:len(base)-1]
	}
	// Find the last segment after '/'
	slashIndex := -1
	for i := len(base) - 1; i >= 0; i-- {
		if base[i] == '/' {
			slashIndex = i
			break
		}
	}
	if slashIndex != -1 && slashIndex < len(base)-1 {
		base = base[slashIndex+1:]
	}
	// Remove .git suffix if present
	if len(base) > 4 && base[len(base)-4:] == ".git" {
		base = base[:len(base)-4]
	}
	return base
}

func CommitAndPushChanges(repoPath string, message string) error {
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
	err = repo.Push(&git.PushOptions{})
	if err != nil {
		return err
	}
	return nil
}
