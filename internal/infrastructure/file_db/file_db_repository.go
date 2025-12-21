package file_db

import (
	"duh/internal/domain/entity"
	"duh/internal/infrastructure/file_db/toml_repo"
	"path/filepath"
	"slices"
)

//TODO:
// - add tests for the GetEnabledRepositories and GetDefaultRepository using the init_db_service to create the tempdir
// - impl other methods of the DbRepository interface

// 	/// Add or update a repository
// 	UpsertRepository(repo entity.Repository) error

// 	/// List all repositories
// 	GetAllRepositories() ([]entity.Repository, error)

// 	/// Delete a repository
// 	DeleteRepository(repoName string) error

// 	/// Rename a repository
// 	// RenameRepository(oldName, newName string) error

// 	/// Set a repository as the default one
// 	ChangeDefaultRepository(repoName string) error

// 	/// Disable a repository from being used
// 	DisableRepository(repoName string) error

// 	/// Enable a repository to be used
// 	EnableRepository(repoName string) error

// 	/// Initialiaze the database if needed
// 	CheckInit() error

type FileDbRepository struct {
	directoryService DirectoryService
	pathProvider     PathProvider
}

func NewFileDbRepository(
	pathProvider PathProvider,
) *FileDbRepository {
	return &FileDbRepository{
		directoryService: *NewDirectoryService(pathProvider),
		pathProvider:     pathProvider,
	}
}

// Implementations of repository.DbRepository

func (f *FileDbRepository) GetEnabledRepositories() ([]entity.Repository, error) {
	allRepoNames, err := f.directoryService.ListRepositoryNames()
	if err != nil {
		return nil, err
	}

	userPrefs, err := f.getUserPreferences()
	if err != nil {
		return nil, err
	}

	enabledRepos := []entity.Repository{}
	for _, repoName := range allRepoNames {
		if !slices.Contains(userPrefs.GetActivatedRepositories(), repoName) {
			continue
		}
		repo, err := f.getRepositoryByName(repoName)
		if err != nil {
			return nil, err
		}
		enabledRepos = append(enabledRepos, *repo)
	}
	return enabledRepos, nil
}

func (f *FileDbRepository) GetDefaultRepository() (*entity.Repository, error) {
	userPrefs, err := f.getUserPreferences()
	if err != nil {
		return nil, err
	}
	defaultRepoName := userPrefs.GetDefaultRepositoryName()
	return f.getRepositoryByName(defaultRepoName)
}

// Helper functions

func getRepoPath(f *FileDbRepository, name string) (string, error) {
	basePath, err := f.getBasePath()
	if err != nil {
		return "", err
	}
	repoDbFilePath := filepath.Join(basePath, "repositories", name, "db.toml")
	return repoDbFilePath, nil
}

func (f *FileDbRepository) getRepositoryByName(name string) (*entity.Repository, error) {
	repoPath, err := getRepoPath(f, name)
	if err != nil {
		return nil, err
	}
	repoToml, err := toml_repo.LoadRepository(repoPath)
	if err != nil {
		return nil, err
	}
	repo := entity.Repository{
		Name:    name,
		Aliases: repoToml.Aliases,
		Exports: repoToml.Exports,
	}
	return &repo, nil
}

func (f *FileDbRepository) getUserPrefPath() (string, error) {
	basePath, err := f.getBasePath()
	if err != nil {
		return "", err
	}
	return filepath.Join(basePath, "user_preferences.toml"), nil
}

func (f *FileDbRepository) getUserPreferences() (*toml_repo.UserPreferenceToml, error) {
	userPrefPath, err := f.getUserPrefPath()
	if err != nil {
		return nil, err
	}
	return toml_repo.LoadUserPreferences(userPrefPath)
}

func (f *FileDbRepository) getBasePath() (string, error) {
	return f.pathProvider.GetPath()
}
