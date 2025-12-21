package integration

import (
	"duh/internal/infrastructure/file_db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) *file_db.FileDbRepository {
	tempdir := t.TempDir() //"tempdir_file_db_repo_tests"
	// os.Mkdir(tempdir, 0755)
	pathProvider := file_db.NewCustomPathProvider(tempdir)
	initService := file_db.NewInitDbService(pathProvider)
	assert.NoError(t, initService.Run())
	return file_db.NewFileDbRepository(pathProvider)
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
