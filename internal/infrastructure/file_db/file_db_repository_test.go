package file_db

import (
	"duh/internal/domain/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) *FileDbRepository {
	tempdir := t.TempDir()
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

func Test_CheckInit(t *testing.T) {
	fileDbRepository := setup(t)
	hasChanged, err := fileDbRepository.CheckInit()
	assert.NoError(t, err)
	assert.Falsef(t, hasChanged, "CheckInit should not make changes on an already initialized DB")
}
