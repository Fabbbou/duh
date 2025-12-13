package file_db

import (
	"duh/internal/domain/entity"
	"duh/internal/infrastructure/file_db/file_dto"
	"duh/internal/infrastructure/file_db/toml_repo"
	"path/filepath"
	"slices"
)

TODO: add tests for the GetEnabledRepositories and GetDefaultRepository using the init_db_service to create the tempdir


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
		if !slices.Contains(userPrefs.ActivatedRepositories, repoName) {
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
	defaultRepoName := userPrefs.DefaultRepositoryName
	return f.getRepositoryByName(defaultRepoName)
}

// Helper functions

func getRepoRepoByName(f *FileDbRepository, name string) (*toml_repo.TomlRepositoryRepository, error) {
	basePath, err := f.getBasePath()
	if err != nil {
		return nil, err
	}
	repoDbFilePath := filepath.Join(basePath, "repositories", name, "db.toml")
	return toml_repo.NewTomlRepositoryRepository(repoDbFilePath), nil
}

func (f *FileDbRepository) getRepositoryByName(name string) (*entity.Repository, error) {
	repoRepo, err := getRepoRepoByName(f, name)
	if err != nil {
		return nil, err
	}
	repo, err := repoRepo.Get()
	return &repo, err
}

func (f *FileDbRepository) getUserPrefRepo() (*toml_repo.TomlUserPreferencesRepository, error) {
	basePath, err := f.getBasePath()
	if err != nil {
		return nil, err
	}
	userPrefFilePath := filepath.Join(basePath, "user_preferences.toml")
	return toml_repo.NewTomlUserPreferencesRepository(userPrefFilePath), nil
}

func (f *FileDbRepository) getUserPreferences() (*file_dto.UserPreferences, error) {
	repo, err := f.getUserPrefRepo()
	if err != nil {
		return nil, err
	}
	prefs, err := repo.Get()
	return &prefs, err
}

func (f *FileDbRepository) getBasePath() (string, error) {
	return f.pathProvider.GetPath()
}
