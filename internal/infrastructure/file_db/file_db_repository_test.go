package file_db

import (
	"duh/internal/domain/entity"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) *FileDbRepository {
	tempdir := filepath.Join(t.TempDir(), "filedbrepo_test")
	// defer os.RemoveAll(tempdir)
	pathProvider := NewCustomPathProvider(tempdir)
	initService := NewInitDbService(pathProvider)
	hasChanged, err := initService.Check()
	assert.NoError(t, err)
	assert.Truef(t, hasChanged, "initialization should have made changes")
	return NewFileDbRepository(pathProvider)
}

func Test_GetEnabledRepositories(t *testing.T) {
	fileDbRepository := setup(t)

	enabledRepos, err := fileDbRepository.GetEnabledRepositories()
	assert.NoError(t, err)

	assert.Lenf(t, enabledRepos, 1, "should get 1 enabled repos (local)")
	assert.Equal(t, "local", enabledRepos[0].Name)
}

func Test_GetDefaultRepository(t *testing.T) {
	fileDbRepository := setup(t)
	defaultRepo, err := fileDbRepository.GetDefaultRepository()
	assert.NoError(t, err)

	assert.Equal(t, "local", defaultRepo.Name)
}

func Test_GetAllRepositories(t *testing.T) {

	fileDbRepository := setup(t)
	fileDbRepository.directoryService.CreateRepository("local2")

	allRepos, err := fileDbRepository.GetAllRepositories()
	assert.NoError(t, err)
	assert.Lenf(t, allRepos, 2, "should get 2 repos (local and local2)")
	assert.Equal(t, "local", allRepos[0].Name)
	assert.Equal(t, "local2", allRepos[1].Name)
}

func Test_DeleteRepository(t *testing.T) {
	fileDbRepository := setup(t)
	repoName := "tobedeleted"
	_, err := fileDbRepository.directoryService.CreateRepository(repoName)
	assert.NoError(t, err)

	//get repo to ensure it exists
	repo, err := fileDbRepository.getRepositoryByName(repoName)
	assert.NoError(t, err)
	assert.Equal(t, repoName, repo.Name)

	//delete repo
	err = fileDbRepository.DeleteRepository(repoName)
	assert.NoError(t, err)

	//get repo again to ensure it no longer exists
	_, err = fileDbRepository.getRepositoryByName(repoName)
	assert.Error(t, err)
}

func Test_UpsertRepository(t *testing.T) {
	fileDbRepository := setup(t)
	repo := entity.Repository{
		Name:    "newrepo",
		Aliases: map[string]string{"nr": "newr"},
		Exports: map[string]string{"export1": "value1"},
	}
	err := fileDbRepository.UpsertRepository(repo)
	assert.NoError(t, err)
	repoP, err := fileDbRepository.getRepositoryByName("newrepo")
	assert.NoError(t, err)
	assert.Equal(t, repo.Name, repoP.Name)
	assert.Equal(t, repo.Aliases, repoP.Aliases)
	assert.Equal(t, repo.Exports, repoP.Exports)
	repoOverride := entity.Repository{
		Name:    "newrepo",
		Aliases: map[string]string{"nr": "newr2"},
		Exports: map[string]string{"export1": "value2"},
	}
	err = fileDbRepository.UpsertRepository(repoOverride)
	assert.NoError(t, err)
	repoP, err = fileDbRepository.getRepositoryByName("newrepo")
	assert.NoError(t, err)
	assert.Equal(t, repoOverride.Name, repoP.Name)
	assert.Equal(t, repoOverride.Aliases, repoP.Aliases)
	assert.Equal(t, repoOverride.Exports, repoP.Exports)
}

func Test_ChangeDefaultRepository(t *testing.T) {
	fileDbRepository := setup(t)
	repoName := "newdefaultrepo"
	_, err := fileDbRepository.directoryService.CreateRepository(repoName)
	assert.NoError(t, err)

	err = fileDbRepository.ChangeDefaultRepository(repoName)
	assert.NoError(t, err)
	defaultRepo, err := fileDbRepository.GetDefaultRepository()
	assert.NoError(t, err)
	assert.Equal(t, repoName, defaultRepo.Name)
}

func Test_EnableRepository(t *testing.T) {
	fileDbRepository := setup(t)

	repoName := "enablerepo"
	_, err := fileDbRepository.directoryService.CreateRepository(repoName)
	assert.NoError(t, err)
	err = fileDbRepository.EnableRepository(repoName)
	assert.NoError(t, err)
	enabledRepos, err := fileDbRepository.GetEnabledRepositories()
	assert.NoError(t, err)
	var found bool
	for _, repo := range enabledRepos {
		if repo.Name == repoName {
			found = true
			break
		}
	}
	assert.Truef(t, found, "enabled repositories should contain the enabled repo")
}

func Test_DisableRepository(t *testing.T) {
	fileDbRepository := setup(t)
	repoName := "disablerepo"
	_, err := fileDbRepository.directoryService.CreateRepository(repoName)
	assert.NoError(t, err)
	err = fileDbRepository.DisableRepository(repoName)
	assert.NoError(t, err)
	enabledRepos, err := fileDbRepository.GetEnabledRepositories()
	assert.NoError(t, err)
	var found bool
	for _, repo := range enabledRepos {
		if repo.Name == repoName {
			found = true
			break
		}
	}
	assert.Falsef(t, found, "enabled repositories should not contain the disabled repo")
}

func Test_AddRepository(t *testing.T) {
	fileDbRepository := setup(t)
	repoURL := "https://github.com/Fabbbou/my-duh"
	repoName, err := fileDbRepository.AddRepository(repoURL, nil)
	assert.NoError(t, err)
	assert.Equal(t, "my-duh", repoName)
	repo, err := fileDbRepository.getRepositoryByName(repoName)
	assert.NoError(t, err)
	assert.Equal(t, repoName, repo.Name)
	assert.NotEmpty(t, repo.Aliases["ll"])
}

func Test_UpdateRepositories(t *testing.T) {
	fileDbRepository := setup(t)

	// Test with no repositories having git remotes
	results, err := fileDbRepository.UpdateRepositories(entity.UpdateSafe)
	assert.NoError(t, err)
	assert.Empty(t, results.LocalChangesDetected)
	assert.Empty(t, results.OtherErrors)
}

func Test_UpdateRepositories_WithGitRepositories(t *testing.T) {
	fileDbRepository := setup(t)

	// Add a repository with git remote
	repoURL := "https://github.com/isomorphic-git/test.empty"
	repoName, err := fileDbRepository.AddRepository(repoURL, nil)
	assert.NoError(t, err)

	// Test safe strategy - should succeed when no local changes
	results, err := fileDbRepository.UpdateRepositories(entity.UpdateSafe)
	assert.NoError(t, err)
	assert.Empty(t, results.LocalChangesDetected)
	assert.Empty(t, results.OtherErrors)

	// Test keep strategy
	results, err = fileDbRepository.UpdateRepositories(entity.UpdateKeep)
	assert.NoError(t, err)
	assert.Empty(t, results.LocalChangesDetected)
	assert.Empty(t, results.OtherErrors)

	// Test force strategy
	results, err = fileDbRepository.UpdateRepositories(entity.UpdateForce)
	assert.NoError(t, err)
	assert.Empty(t, results.LocalChangesDetected)
	assert.Empty(t, results.OtherErrors)

	// Cleanup
	fileDbRepository.DeleteRepository(repoName)
}

func Test_UpdateRepositories_InvalidStrategy(t *testing.T) {
	fileDbRepository := setup(t)
	// Add a repository with git remote
	repoURL := "https://github.com/isomorphic-git/test.empty"
	repoName, err := fileDbRepository.AddRepository(repoURL, nil)
	assert.NoError(t, err)

	// Create local changes to trigger strategy validation
	basePath, err := fileDbRepository.pathProvider.GetPath()
	assert.NoError(t, err)
	repoPath := filepath.Join(basePath, "repositories", repoName)
	testFile := filepath.Join(repoPath, "local-change.txt")
	err = os.WriteFile(testFile, []byte("local changes"), 0644)
	assert.NoError(t, err)

	// Test with invalid strategy - should return error
	results, err := fileDbRepository.UpdateRepositories("invalid")
	assert.NoError(t, err) // The function itself doesn't error, but individual repos might
	assert.Empty(t, results.LocalChangesDetected)
	// Should have an error for the repository with invalid strategy
	assert.NotEmpty(t, results.OtherErrors)

	// Cleanup
	fileDbRepository.DeleteRepository(repoName)
}

func Test_PushRepository_NoGitRepo(t *testing.T) {
	fileDbRepository := setup(t)

	// Try to push a repository without git
	err := fileDbRepository.PushRepository("local")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a git repository")
}

func Test_PushRepository_NoRemote(t *testing.T) {
	fileDbRepository := setup(t)

	// Create a local git repository without remote
	repoName := "local-git-only"
	repoPath, err := fileDbRepository.directoryService.CreateRepository(repoName)
	assert.NoError(t, err)

	// Initialize as git repository but don't add remote
	_, err = git.PlainInit(repoPath, false)
	assert.NoError(t, err)

	// Try to push - should fail because no remote
	err = fileDbRepository.PushRepository(repoName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not have a git remote configured")
}

func Test_PushRepository_RepositoryNotFound(t *testing.T) {
	fileDbRepository := setup(t)

	// Try to push non-existent repository
	err := fileDbRepository.PushRepository("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository 'nonexistent' not found")
}
