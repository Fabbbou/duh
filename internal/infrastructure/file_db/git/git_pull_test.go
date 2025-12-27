package gitt

import (
	"duh/internal/domain/entity"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
)

func Test_hasUpdatesPending_NoUpdates(t *testing.T) {
	// Test with a fresh clone - should have no updates pending
	outputPath := t.TempDir()
	defer os.RemoveAll(outputPath)

	repoPath := filepath.Join(outputPath, "test-repo")
	err := CloneGitRepository("https://github.com/isomorphic-git/test.empty", repoPath)
	assert.NoError(t, err)

	// Fresh clone should have no pending updates
	hasUpdates, err := hasUpdatesPending(repoPath)
	assert.NoError(t, err)
	assert.False(t, hasUpdates, "Fresh clone should not have pending updates")
}

func Test_hasUpdatesPending_WithUpdates(t *testing.T) {
	// Create a test scenario where updates are pending
	outputPath := t.TempDir()
	defer os.RemoveAll(outputPath)

	repoPath := filepath.Join(outputPath, "test-repo")

	// Clone a repository that has multiple commits (using a test repo with history)
	err := CloneGitRepository("https://github.com/octocat/Hello-World", repoPath)
	assert.NoError(t, err)

	// Open the repository
	repo, err := git.PlainOpen(repoPath)
	assert.NoError(t, err)

	// Get the current HEAD
	head, err := repo.Head()
	assert.NoError(t, err)

	// Get commit history to find a previous commit
	iter, err := repo.Log(&git.LogOptions{From: head.Hash()})
	assert.NoError(t, err)

	var commits []plumbing.Hash
	err = iter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c.Hash)
		return nil
	})
	assert.NoError(t, err)

	// If we have at least 2 commits, reset to an earlier one
	if len(commits) > 1 {
		// Get worktree and reset to previous commit
		w, err := repo.Worktree()
		assert.NoError(t, err)

		// Reset to the second commit (making local behind remote)
		err = w.Reset(&git.ResetOptions{
			Commit: commits[1],
			Mode:   git.HardReset,
		})
		assert.NoError(t, err)

		// Now test if updates are pending - should return true
		hasUpdates, err := hasUpdatesPending(repoPath)
		assert.NoError(t, err)
		assert.True(t, hasUpdates, "Repository should have pending updates after reset")
	} else {
		t.Skip("Test repository doesn't have enough commit history")
	}
}

// Test helper function to create a test repository with content
func createTestRepoWithContent(t *testing.T, repoPath string) {
	// Initialize a new repository
	repo, err := git.PlainInit(repoPath, false)
	assert.NoError(t, err)

	// Create a test file
	testFile := filepath.Join(repoPath, "test.txt")
	err = os.WriteFile(testFile, []byte("initial content"), 0644)
	assert.NoError(t, err)

	// Add and commit the file
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
}

func Test_PullAllRepositories_NoRepositories(t *testing.T) {
	// Test with empty directory
	tempDir := t.TempDir()

	results, err := PullAllRepositories(tempDir, entity.UpdateSafe)
	assert.NoError(t, err)
	assert.Empty(t, results.LocalChangesDetected)
	assert.Empty(t, results.OtherErrors)
}

func Test_PullAllRepositories_NonGitRepositories(t *testing.T) {
	// Create test directory with non-git repositories
	tempDir := t.TempDir()

	// Create some directories without git
	nonGitDir1 := filepath.Join(tempDir, "not-a-repo")
	nonGitDir2 := filepath.Join(tempDir, "another-folder")

	err := os.Mkdir(nonGitDir1, 0755)
	assert.NoError(t, err)
	err = os.Mkdir(nonGitDir2, 0755)
	assert.NoError(t, err)

	results, err := PullAllRepositories(tempDir, entity.UpdateSafe)
	assert.NoError(t, err)
	assert.Empty(t, results.LocalChangesDetected)
	assert.Empty(t, results.OtherErrors)
}

func Test_PullAllRepositories_GitRepoWithoutRemote(t *testing.T) {
	tempDir := t.TempDir()

	// Create a git repository without remote
	repoPath := filepath.Join(tempDir, "local-repo")
	createTestRepoWithContent(t, repoPath)

	results, err := PullAllRepositories(tempDir, entity.UpdateSafe)
	assert.NoError(t, err)
	assert.Empty(t, results.LocalChangesDetected)
	assert.Empty(t, results.OtherErrors)
}

func Test_PullAllRepositories_WithRemoteRepository(t *testing.T) {
	tempDir := t.TempDir()

	// Clone a repository with remote
	repoPath := filepath.Join(tempDir, "remote-repo")
	err := CloneGitRepository("https://github.com/isomorphic-git/test.empty", repoPath)
	assert.NoError(t, err)

	// Test safe strategy
	results, err := PullAllRepositories(tempDir, entity.UpdateSafe)
	assert.NoError(t, err)
	assert.Empty(t, results.LocalChangesDetected)
	assert.Empty(t, results.OtherErrors)

	// Test keep strategy
	results, err = PullAllRepositories(tempDir, entity.UpdateKeep)
	assert.NoError(t, err)
	assert.Empty(t, results.LocalChangesDetected)
	assert.Empty(t, results.OtherErrors)

	// Test force strategy
	results, err = PullAllRepositories(tempDir, entity.UpdateForce)
	assert.NoError(t, err)
	assert.Empty(t, results.LocalChangesDetected)
	assert.Empty(t, results.OtherErrors)
}

func Test_PullAllRepositories_WithLocalChanges(t *testing.T) {
	tempDir := t.TempDir()

	// Clone a repository
	repoPath := filepath.Join(tempDir, "repo-with-changes")
	err := CloneGitRepository("https://github.com/isomorphic-git/test.empty", repoPath)
	assert.NoError(t, err)

	// Create local changes
	testFile := filepath.Join(repoPath, "local-change.txt")
	err = os.WriteFile(testFile, []byte("local changes"), 0644)
	assert.NoError(t, err)

	// Test safe strategy - should detect local changes
	results, err := PullAllRepositories(tempDir, entity.UpdateSafe)
	assert.NoError(t, err)
	assert.Contains(t, results.LocalChangesDetected, "repo-with-changes")
	assert.Empty(t, results.OtherErrors)

	// Test keep strategy - should commit and pull
	results, err = PullAllRepositories(tempDir, entity.UpdateKeep)
	assert.NoError(t, err)
	assert.Empty(t, results.LocalChangesDetected)
	assert.Empty(t, results.OtherErrors)

	// Verify the file was committed
	repo, err := git.PlainOpen(repoPath)
	assert.NoError(t, err)
	w, err := repo.Worktree()
	assert.NoError(t, err)
	status, err := w.Status()
	assert.NoError(t, err)
	assert.True(t, status.IsClean())
}

func Test_PullAllRepositories_ForceStrategy(t *testing.T) {
	tempDir := t.TempDir()

	// Clone a repository
	repoPath := filepath.Join(tempDir, "repo-force-test")
	err := CloneGitRepository("https://github.com/isomorphic-git/test.empty", repoPath)
	assert.NoError(t, err)

	// Create local changes
	testFile := filepath.Join(repoPath, "local-change.txt")
	err = os.WriteFile(testFile, []byte("local changes"), 0644)
	assert.NoError(t, err)

	// Test force strategy - should discard local changes
	results, err := PullAllRepositories(tempDir, entity.UpdateForce)
	assert.NoError(t, err)
	assert.Empty(t, results.LocalChangesDetected)
	assert.Empty(t, results.OtherErrors)

	// Verify local changes were discarded
	_, err = os.Stat(testFile)
	assert.True(t, os.IsNotExist(err))
}

func Test_PullAllRepositories_MultipleRepositories(t *testing.T) {
	tempDir := t.TempDir()

	// Create multiple repositories
	repo1Path := filepath.Join(tempDir, "repo1")
	repo2Path := filepath.Join(tempDir, "repo2")
	repo3Path := filepath.Join(tempDir, "local-only")

	// Clone remote repositories
	err := CloneGitRepository("https://github.com/isomorphic-git/test.empty", repo1Path)
	assert.NoError(t, err)
	err = CloneGitRepository("https://github.com/isomorphic-git/test.empty", repo2Path)
	assert.NoError(t, err)

	// Create local-only repository
	createTestRepoWithContent(t, repo3Path)

	// Add local changes to repo2
	testFile := filepath.Join(repo2Path, "local-change.txt")
	err = os.WriteFile(testFile, []byte("local changes"), 0644)
	assert.NoError(t, err)

	// Test safe strategy
	results, err := PullAllRepositories(tempDir, entity.UpdateSafe)
	assert.NoError(t, err)
	assert.Contains(t, results.LocalChangesDetected, "repo2")
	assert.NotContains(t, results.LocalChangesDetected, "repo1")
	assert.NotContains(t, results.LocalChangesDetected, "local-only") // Local-only repo shouldn't be processed
	assert.Empty(t, results.OtherErrors)
}

func Test_PullAllRepositories_InvalidStrategy(t *testing.T) {
	tempDir := t.TempDir()

	// Clone a repository
	repoPath := filepath.Join(tempDir, "repo-invalid-strategy")
	err := CloneGitRepository("https://github.com/isomorphic-git/test.empty", repoPath)
	assert.NoError(t, err)

	// Create local changes to trigger strategy validation
	testFile := filepath.Join(repoPath, "local-change.txt")
	err = os.WriteFile(testFile, []byte("local changes"), 0644)
	assert.NoError(t, err)

	// Test with invalid strategy
	results, err := PullAllRepositories(tempDir, "invalid-strategy")
	assert.NoError(t, err)
	assert.Empty(t, results.LocalChangesDetected)
	assert.Len(t, results.OtherErrors, 1)
	assert.Contains(t, results.OtherErrors[0].Error(), "unknown strategy")
}
