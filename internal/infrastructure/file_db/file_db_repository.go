package file_db

import (
	"duh/internal/domain/entity"
	"duh/internal/infrastructure/file_db/toml_repo"
	"os"
	"path/filepath"
	"slices"
)

type FileDbRepository struct {
	directoryService DirectoryService
	pathProvider     PathProvider
}

func NewFileDbRepository(pathProvider PathProvider) *FileDbRepository {
	return &FileDbRepository{
		directoryService: *NewDirectoryService(pathProvider),
		pathProvider:     pathProvider,
	}
}

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

// Add or update a repository
func (f *FileDbRepository) UpsertRepository(repo entity.Repository) error {
	repoToml := toml_repo.RepositoryToml{
		Aliases: repo.Aliases,
		Exports: repo.Exports,
	}

	repoPath, err := f.directoryService.CreateRepository(repo.Name)
	if err != nil {
		return err
	}
	dbPath := filepath.Join(repoPath, "db.toml")
	return toml_repo.SaveToml(dbPath, &repoToml)
}

func (f *FileDbRepository) DeleteRepository(repoName string) error {
	return f.directoryService.DeleteRepository(repoName)
}

// Set a repository as the default one
func (f *FileDbRepository) ChangeDefaultRepository(repoName string) error {
	userPrefs, err := f.getUserPreferences()
	if err != nil {
		return err
	}
	userPrefs.SetDefaultRepositoryName(repoName)
	userPrefPath, err := f.getUserPrefPath()
	if err != nil {
		return err
	}
	return toml_repo.SaveToml(userPrefPath, userPrefs)
}

// Enable a repository to be used
func (f *FileDbRepository) EnableRepository(repoName string) error {
	userPrefs, err := f.getUserPreferences()
	if err != nil {
		return err
	}
	activatedRepos := userPrefs.GetActivatedRepositories()
	if !slices.Contains(activatedRepos, repoName) {
		activatedRepos = append(activatedRepos, repoName)
		userPrefs.SetActivatedRepositories(activatedRepos)
	}
	userPrefPath, err := f.getUserPrefPath()
	if err != nil {
		return err
	}
	return toml_repo.SaveToml(userPrefPath, userPrefs)
}

// Disable a repository from being used
func (f *FileDbRepository) DisableRepository(repoName string) error {
	userPrefs, err := f.getUserPreferences()
	if err != nil {
		return err
	}
	activatedRepos := userPrefs.GetActivatedRepositories()
	if slices.Contains(activatedRepos, repoName) {
		newActivatedRepos := []string{}
		for _, r := range activatedRepos {
			if r != repoName {
				newActivatedRepos = append(newActivatedRepos, r)
			}
		}
		userPrefs.SetActivatedRepositories(newActivatedRepos)
	}
	userPrefPath, err := f.getUserPrefPath()
	if err != nil {
		return err
	}
	return toml_repo.SaveToml(userPrefPath, userPrefs)
}

// Initialiaze the database if needed
func (f *FileDbRepository) CheckInit() (bool, error) {
	initService := NewInitDbService(f.pathProvider)
	return initService.Check()
}

// Rename a repository
func (f *FileDbRepository) RenameRepository(oldName, newName string) error {
	oldRepoPath, err := f.directoryService.getRepositoryPath(oldName)
	if err != nil {
		return err
	}
	newRepoPath, err := f.directoryService.getRepositoryPath(newName)
	if err != nil {
		return err
	}

	return os.Rename(oldRepoPath, newRepoPath)
}

////////////////////
// Helper (internal) functions
////////////////////

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

func (f *FileDbRepository) GetAllRepositories() ([]entity.Repository, error) {
	allRepoNames, err := f.directoryService.ListRepositoryNames()
	if err != nil {
		return nil, err
	}
	allRepos := []entity.Repository{}
	for _, repoName := range allRepoNames {
		repo, err := f.getRepositoryByName(repoName)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, *repo)
	}
	return allRepos, nil
}
