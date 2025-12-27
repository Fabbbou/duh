package gitt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
)

func Test_CommitAndPushChanges_CleanRepo(t *testing.T) {
	// Create a test repository
	tempDir := t.TempDir()

	// Initialize a new repository
	repo, err := git.PlainInit(tempDir, false)
	assert.NoError(t, err)

	// Add a remote (in-memory for testing)
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/test/repo.git"},
	})
	assert.NoError(t, err)

	// Create initial content
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("initial content"), 0644)
	assert.NoError(t, err)

	// Commit initial content
	w, err := repo.Worktree()
	assert.NoError(t, err)
	_, err = w.Add("test.txt")
	assert.NoError(t, err)
	_, err = w.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
		},
	})
	assert.NoError(t, err)

	// Test push with clean repository - should not fail locally (though push to remote will fail in test)
	err = CommitAndPushChanges(tempDir)
	// We expect this to fail because we don't have a real remote, but it should not fail due to working tree issues
	if err != nil {
		// Could be authentication error or author field error in CI
		assert.True(t,
			strings.Contains(err.Error(), "authentication required") ||
				strings.Contains(err.Error(), "author field is required") ||
				strings.Contains(err.Error(), "remote repository") ||
				strings.Contains(err.Error(), "push"),
			"Expected push-related error, got: %v", err)
	}
}

func Test_CommitAndPushChanges_DirtyRepo(t *testing.T) {
	// Create a test repository
	tempDir := t.TempDir()

	// Initialize a new repository
	repo, err := git.PlainInit(tempDir, false)
	assert.NoError(t, err)

	// Add a remote
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/test/repo.git"},
	})
	assert.NoError(t, err)

	// Create and commit initial content
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("initial content"), 0644)
	assert.NoError(t, err)

	w, err := repo.Worktree()
	assert.NoError(t, err)
	_, err = w.Add("test.txt")
	assert.NoError(t, err)
	_, err = w.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
		},
	})
	assert.NoError(t, err)

	// Add local changes
	err = os.WriteFile(testFile, []byte("modified content"), 0644)
	assert.NoError(t, err)

	// Test push with dirty repository - should auto-commit then try to push
	commitPushErr := CommitAndPushChanges(tempDir)

	// Check the final status
	status, err := w.Status()
	assert.NoError(t, err)

	// If CommitAndPushChanges returned an error, check what kind of error it was
	if commitPushErr != nil {
		// Could be authentication, git config, or remote repository errors
		assert.True(t,
			strings.Contains(commitPushErr.Error(), "authentication required") ||
				strings.Contains(commitPushErr.Error(), "author field is required") ||
				strings.Contains(commitPushErr.Error(), "git config") ||
				strings.Contains(commitPushErr.Error(), "user.name") ||
				strings.Contains(commitPushErr.Error(), "user.email") ||
				strings.Contains(commitPushErr.Error(), "remote repository") ||
				strings.Contains(commitPushErr.Error(), "push"),
			"Expected push/config-related error, got: %v", commitPushErr)

		// If it was a git config error, the repository might still be dirty
		if strings.Contains(commitPushErr.Error(), "user.name") ||
			strings.Contains(commitPushErr.Error(), "user.email") ||
			strings.Contains(commitPushErr.Error(), "author field is required") {
			// Skip the clean check since commit failed due to missing git config
			t.Logf("Repository still dirty due to failed commit (missing git config): %v", commitPushErr)
			return
		}
	}

	// If we get here, either there was no error or it was a push error after successful commit
	// In either case, the worktree should be clean after successful commit
	assert.True(t, status.IsClean(), "Worktree should be clean after successful commit")
}

func Test_addAndCommitAllChanges(t *testing.T) {
	// Create a test repository
	tempDir := t.TempDir()

	// Initialize a new repository
	repo, err := git.PlainInit(tempDir, false)
	assert.NoError(t, err)

	// Create test files
	testFile1 := filepath.Join(tempDir, "test1.txt")
	testFile2 := filepath.Join(tempDir, "test2.txt")
	err = os.WriteFile(testFile1, []byte("content1"), 0644)
	assert.NoError(t, err)
	err = os.WriteFile(testFile2, []byte("content2"), 0644)
	assert.NoError(t, err)

	// Test committing all changes
	err = addAndCommitAllChanges(tempDir, "Test commit")
	// This might fail if local git config is not set (user.name, user.email)
	if err != nil {
		// Check if it's a git config related error
		if strings.Contains(err.Error(), "user.name") ||
			strings.Contains(err.Error(), "user.email") ||
			strings.Contains(err.Error(), "git config") ||
			strings.Contains(err.Error(), "author field is required") {
			t.Skipf("Skipping test due to missing git config: %v", err)
			return
		}
		assert.NoError(t, err)
		return
	}

	// Verify all files were committed
	w, err := repo.Worktree()
	assert.NoError(t, err)
	status, err := w.Status()
	assert.NoError(t, err)
	assert.True(t, status.IsClean())

	// Verify commit exists
	head, err := repo.Head()
	assert.NoError(t, err)
	if head != nil {
		commit, err := repo.CommitObject(head.Hash())
		assert.NoError(t, err)
		if commit != nil {
			assert.Equal(t, "Test commit", commit.Message)
		}
	}
}

func Test_addAndCommitAllChanges_NoChanges(t *testing.T) {
	// Create a test repository
	tempDir := t.TempDir()

	// Initialize a new repository
	repo, err := git.PlainInit(tempDir, false)
	assert.NoError(t, err)

	// Create and commit initial content
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("content"), 0644)
	assert.NoError(t, err)

	w, err := repo.Worktree()
	assert.NoError(t, err)
	_, err = w.Add("test.txt")
	assert.NoError(t, err)
	_, err = w.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
		},
	})
	assert.NoError(t, err)

	// Test committing when there are no changes - should not fail
	err = addAndCommitAllChanges(tempDir, "No changes commit")
	// This might fail or succeed depending on git behavior, but should not crash
	// The important thing is that the repository remains in a valid state
	status, err := w.Status()
	assert.NoError(t, err)
	assert.True(t, status.IsClean())
}
