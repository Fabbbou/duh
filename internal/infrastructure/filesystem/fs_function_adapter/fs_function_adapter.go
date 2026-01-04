package fs_function_adapter

import (
	"duh/internal/domain/entity"
	"duh/internal/infrastructure/filesystem/common"
	"duh/internal/infrastructure/filesystem/fs_user_repository"
	"duh/internal/infrastructure/filesystem/function"
	"fmt"
	"path/filepath"
)

type FSFunctionAdapter struct {
	pathProvider             common.PathProvider
	userPreferenceRepository *fs_user_repository.FsUserRepository
	directoryService         *common.DirectoryService
}

func NewFSFunctionsRepository(
	pathProvider common.PathProvider,
	userPreferenceRepository *fs_user_repository.FsUserRepository,
) *FSFunctionAdapter {
	return &FSFunctionAdapter{
		pathProvider:             pathProvider,
		userPreferenceRepository: userPreferenceRepository,
		directoryService:         common.NewDirectoryService(pathProvider),
	}
}

func GetInternalScripts() ([]entity.Script, error) {
	scriptRequire, err := function.GetScriptFromString("require", function.RequireShScript, "duh://internal-script/require.sh")
	if err != nil {
		return nil, err
	}
	return []entity.Script{*scriptRequire}, nil
}

func (f *FSFunctionAdapter) GetActivatedScripts() ([]entity.Script, error) {
	userPrefs, err := f.userPreferenceRepository.GetUserPreference()
	if err != nil {
		return nil, err
	}
	return f.getScriptsForRepos(userPrefs.Repositories.ActivatedRepositories)
}

func (f *FSFunctionAdapter) GetAllScripts() ([]entity.Script, error) {
	repoDirs, err := f.directoryService.ListRepositoryNames()
	if err != nil {
		return nil, err
	}
	return f.getScriptsForRepos(repoDirs)
}

func (f *FSFunctionAdapter) GetFunctionsPath(repoName string) (string, error) {
	path, err := f.pathProvider.GetPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(path, "repositories", repoName, "functions"), nil
}

func (f *FSFunctionAdapter) getScriptsForRepos(repoNames []string) ([]entity.Script, error) {
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

func (f *FSFunctionAdapter) GetInternalScripts() ([]entity.Script, error) {
	return GetInternalScripts()
}
