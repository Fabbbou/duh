package fs_functions_repository

import (
	"duh/internal/domain/entity"
	"duh/internal/infrastructure/filesystem/common"
	"duh/internal/infrastructure/filesystem/fs_user_repository"
	"duh/internal/infrastructure/filesystem/function"
	"path/filepath"
)

type FSFunctionsRepository struct {
	pathProvider             common.PathProvider
	userPreferenceRepository *fs_user_repository.FsUserRepository
}

func NewFSFunctionsRepository(
	pathProvider common.PathProvider,
	userPreferenceRepository *fs_user_repository.FsUserRepository,
) *FSFunctionsRepository {
	return &FSFunctionsRepository{
		pathProvider:             pathProvider,
		userPreferenceRepository: userPreferenceRepository,
	}
}

func (f *FSFunctionsRepository) GetActivatedScripts() ([]entity.Script, error) {
	userPrefs, err := f.userPreferenceRepository.GetUserPreference()
	if err != nil {
		return nil, err
	}
	path, err := f.pathProvider.GetPath()
	if err != nil {
		return nil, err
	}
	repoPath := filepath.Join(path, "repositories")
	var scripts []entity.Script
	for _, repoName := range userPrefs.Repositories.ActivatedRepositories {
		scriptDirPath := filepath.Join(repoPath, repoName, "functions")
		script, err := function.GetScripts(scriptDirPath)
		if err != nil {
			continue
		}
		scripts = append(scripts, script...)
	}
	return scripts, nil
}
