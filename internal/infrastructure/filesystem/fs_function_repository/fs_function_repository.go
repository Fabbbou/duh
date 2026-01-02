package fs_functions_repository

import (
	"duh/internal/domain/entity"
	"duh/internal/infrastructure/filesystem/common"
	"duh/internal/infrastructure/filesystem/fs_user_repository"
	"duh/internal/infrastructure/filesystem/function"
	"fmt"
	"path/filepath"
)

type FSFunctionsRepository struct {
	pathProvider             common.PathProvider
	userPreferenceRepository *fs_user_repository.FsUserRepository
	directoryService         *common.DirectoryService
}

func NewFSFunctionsRepository(
	pathProvider common.PathProvider,
	userPreferenceRepository *fs_user_repository.FsUserRepository,
) *FSFunctionsRepository {
	return &FSFunctionsRepository{
		pathProvider:             pathProvider,
		userPreferenceRepository: userPreferenceRepository,
		directoryService:         common.NewDirectoryService(pathProvider),
	}
}

func (f *FSFunctionsRepository) GetActivatedScripts() ([]entity.Script, error) {
	userPrefs, err := f.userPreferenceRepository.GetUserPreference()
	if err != nil {
		return nil, err
	}
	return f.getScriptsForRepos(userPrefs.Repositories.ActivatedRepositories)
}

func (f *FSFunctionsRepository) GetAllScripts() ([]entity.Script, error) {
	repoDirs, err := f.directoryService.ListRepositoryNames()
	if err != nil {
		return nil, err
	}
	return f.getScriptsForRepos(repoDirs)
}

func (f *FSFunctionsRepository) GetFunctionsPath(repoName string) (string, error) {
	path, err := f.pathProvider.GetPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(path, "repositories", repoName, "functions"), nil
}

func (f *FSFunctionsRepository) getScriptsForRepos(repoNames []string) ([]entity.Script, error) {
	var scripts []entity.Script
	var errors []error
	for _, repoName := range repoNames {
		scriptDirPath, err := f.GetFunctionsPath(repoName)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		script, err := function.GetScripts(scriptDirPath)
		if err == function.ErrDirNotFound {
			// Skip missing directories
			continue
		}
		if err != nil {
			errors = append(errors, err)
			continue
		}
		scripts = append(scripts, script...)
	}
	if len(errors) > 0 {
		return scripts, fmt.Errorf("errors occurred while getting activated scripts: %v", errors)
	}
	return scripts, nil
}
