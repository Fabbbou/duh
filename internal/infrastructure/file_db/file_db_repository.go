package file_db

import (
	"duh/internal/domain/entity"
	"duh/internal/infrastructure/editor"
	gitt "duh/internal/infrastructure/file_db/git"
	"duh/internal/infrastructure/file_db/toml_repo"
	"fmt"
	"os"
	"os/exec"
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

func (f *FileDbRepository) AddRepository(url string, name *string) (string, error) {
	path, err := f.getBasePath()
	if err != nil {
		return "", err
	}
	finalName := ""
	if name == nil || len(*name) <= 0 {
		finalName = gitt.ExtractGitRepoName(url)
	} else {
		finalName = *name
	}
	if finalName == "" {
		return "", fmt.Errorf("cannot add a repo without a name")
	}
	repoPath := filepath.Join(path, "repositories", finalName)
	return finalName, gitt.CloneGitRepository(url, repoPath)
}

func (f *FileDbRepository) CreateRepository(name string) (string, error) {
	repo, err := f.getRepositoryByName(name)
	if err == nil && repo != nil {
		return "", fmt.Errorf("repository with name '%s' already exists", name)
	}

	repoPath, err := f.directoryService.CreateRepository(name)
	if err != nil {
		return "", err
	}
	// Initialize empty toml file
	repoToml := toml_repo.RepositoryToml{
		Aliases: map[string]string{},
		Exports: map[string]string{},
		Metadata: toml_repo.MetadataMap{
			NameOrigin: name,
		},
	}
	dbPath := filepath.Join(repoPath, "db.toml")
	err = toml_repo.SaveToml(dbPath, &repoToml)
	if err != nil {
		return "", err
	}

	return repoPath, f.EnableRepository(name)
}

func (f *FileDbRepository) UpdateRepositories(strategy string) (entity.RepositoryUpdateResults, error) {
	path, err := f.getBasePath()
	if err != nil {
		return entity.RepositoryUpdateResults{}, err
	}
	reposPath := filepath.Join(path, "repositories")

	return gitt.PullAllRepositories(reposPath, strategy)
}

func (f *FileDbRepository) EditRepo(repoName string) error {
	// Check if repository exists
	_, err := f.getRepositoryByName(repoName)
	if err != nil {
		return fmt.Errorf("repository '%s' not found: %w", repoName, err)
	}

	// Get repository db.toml file path
	repoDbFilePath, err := createRepoDbFilePath(f, repoName)
	if err != nil {
		return fmt.Errorf("failed to get repository path: %w", err)
	}

	// Find default editor
	editorCmd := editor.FindDefaultFileEditor()

	// Create and run editor command
	cmd := exec.Command(editorCmd, repoDbFilePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to open editor '%s': %w", editorCmd, err)
	}

	return nil
}

////////////////////
// Helper (internal) functions
////////////////////

func createRepoDbFilePath(f *FileDbRepository, name string) (string, error) {
	basePath, err := f.getBasePath()
	if err != nil {
		return "", err
	}
	repoDbFilePath := filepath.Join(basePath, "repositories", name, "db.toml")
	return repoDbFilePath, nil
}

func (f *FileDbRepository) getRepositoryByName(name string) (*entity.Repository, error) {
	repoPath, err := createRepoDbFilePath(f, name)
	if err != nil {
		return nil, err
	}
	repoToml, err := toml_repo.LoadRepository(repoPath)
	if err != nil {
		return nil, err
	}
	aliases := map[string]string{}
	if repoToml.Aliases == nil {
		aliases = map[string]string{}
	} else {
		aliases = repoToml.Aliases
	}
	exports := map[string]string{}
	if repoToml.Exports == nil {
		exports = map[string]string{}
	} else {
		exports = repoToml.Exports
	}
	repo := entity.Repository{
		Name:    name,
		Aliases: aliases,
		Exports: exports,
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
