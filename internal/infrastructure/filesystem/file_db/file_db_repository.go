package file_db

import (
	"duh/internal/domain/entity"
	"duh/internal/domain/utils/gitconfig"
	"duh/internal/infrastructure/filesystem/common"
	"duh/internal/infrastructure/filesystem/editor"
	gitt "duh/internal/infrastructure/filesystem/gitt"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"

	"github.com/go-git/go-git/v5"
)

type FileDbRepository struct {
	DirectoryService      common.DirectoryService
	PathProvider          common.PathProvider
	gitConfigPathProvider common.PathProvider
	fileHandler           FileHandler
}

func NewFileDbRepository(
	PathProvider common.PathProvider,
	gitConfigPathProvider common.PathProvider,
	fileHandler FileHandler,
) *FileDbRepository {
	return &FileDbRepository{
		DirectoryService:      *common.NewDirectoryService(PathProvider),
		PathProvider:          PathProvider,
		gitConfigPathProvider: gitConfigPathProvider,
		fileHandler:           fileHandler,
	}
}

func (f *FileDbRepository) GetEnabledRepositories() ([]entity.Repository, error) {
	allRepoNames, err := f.DirectoryService.ListRepositoryNames()
	if err != nil {
		return nil, err
	}

	userPrefs, err := f.getUserPreferences()
	if err != nil {
		return nil, err
	}

	enabledRepos := []entity.Repository{}
	for _, repoName := range allRepoNames {
		if !slices.Contains(userPrefs.Repositories.ActivatedRepositories, repoName) {
			continue
		}
		repo, err := f.GetRepositoryByName(repoName)
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
	defaultRepoName := userPrefs.Repositories.DefaultRepositoryName
	return f.GetRepositoryByName(defaultRepoName)
}

// Add or update a repository
func (f *FileDbRepository) UpsertRepository(repo entity.Repository) error {
	repoDto := RepositoryDto{
		Aliases: repo.Aliases,
		Exports: repo.Exports,
	}

	repoPath, err := f.DirectoryService.CreateRepository(repo.Name)
	if err != nil {
		return err
	}
	file_name := "db." + f.fileHandler.Extension()
	dbPath := filepath.Join(repoPath, file_name)
	return f.fileHandler.SaveRepositoryFile(dbPath, &repoDto)
}

func (f *FileDbRepository) DeleteRepository(repoName string) error {
	return f.DirectoryService.DeleteRepository(repoName)
}

// Set a repository as the default one
func (f *FileDbRepository) ChangeDefaultRepository(repoName string) error {
	userPrefs, err := f.getUserPreferences()
	if err != nil {
		return err
	}
	userPrefs.Repositories.DefaultRepositoryName = repoName
	userPrefPath, err := f.getUserPrefPath()
	if err != nil {
		return err
	}
	return f.fileHandler.SaveUserPreferenceFile(userPrefPath, userPrefs)
}

// Enable a repository to be used
func (f *FileDbRepository) EnableRepository(repoName string) error {
	userPrefs, err := f.getUserPreferences()
	if err != nil {
		return err
	}
	activatedRepos := userPrefs.Repositories.ActivatedRepositories
	if !slices.Contains(activatedRepos, repoName) {
		activatedRepos = append(activatedRepos, repoName)
		userPrefs.Repositories.ActivatedRepositories = activatedRepos
	}
	userPrefPath, err := f.getUserPrefPath()
	if err != nil {
		return err
	}
	return f.fileHandler.SaveUserPreferenceFile(userPrefPath, userPrefs)
}

// Disable a repository from being used
func (f *FileDbRepository) DisableRepository(repoName string) error {
	userPrefs, err := f.getUserPreferences()
	if err != nil {
		return err
	}
	activatedRepos := userPrefs.Repositories.ActivatedRepositories
	if slices.Contains(activatedRepos, repoName) {
		newActivatedRepos := []string{}
		for _, r := range activatedRepos {
			if r != repoName {
				newActivatedRepos = append(newActivatedRepos, r)
			}
		}
		userPrefs.Repositories.ActivatedRepositories = newActivatedRepos
	}
	userPrefPath, err := f.getUserPrefPath()
	if err != nil {
		return err
	}
	return f.fileHandler.SaveUserPreferenceFile(userPrefPath, userPrefs)
}

// Rename a repository
func (f *FileDbRepository) RenameRepository(oldName, newName string) error {
	oldRepoPath, err := f.DirectoryService.GetRepositoryPath(oldName)
	if err != nil {
		return err
	}
	newRepoPath, err := f.DirectoryService.GetRepositoryPath(newName)
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
	repo, err := f.GetRepositoryByName(name)
	if err == nil && repo != nil {
		return "", fmt.Errorf("repository with name '%s' already exists", name)
	}

	repoPath, err := f.DirectoryService.CreateRepository(name)
	if err != nil {
		return "", err
	}
	// Initialize empty file
	repoDto := RepositoryDto{
		Aliases: map[string]string{},
		Exports: map[string]string{},
		Metadata: MetadataDto{
			NameOrigin: name,
		},
	}
	fileName := "db." + f.fileHandler.Extension()
	dbPath := filepath.Join(repoPath, fileName)
	err = f.fileHandler.SaveRepositoryFile(dbPath, &repoDto)
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

func (f *FileDbRepository) editFile(filePath string) error {
	// Find default editor
	editorCmd := editor.FindDefaultFileEditor()

	// Create and run editor command
	cmd := exec.Command(editorCmd, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to open editor '%s': %w", editorCmd, err)
	}

	return nil
}

func (f *FileDbRepository) EditGitconfig(repoName string) error {
	// Check if repository exists
	_, err := f.GetRepositoryByName(repoName)
	if err != nil {
		return fmt.Errorf("repository '%s' not found: %w", repoName, err)
	}

	gitconfigFile := f.getRepositoryGitconfigPath(repoName)
	if gitconfigFile == "" {
		err = f.DirectoryService.CreateGitconfigFile(repoName)
		if err != nil {
			return fmt.Errorf("failed to create gitconfig file for repository '%s': %w", repoName, err)
		}
		gitconfigFile = f.getRepositoryGitconfigPath(repoName)
	}

	return f.editFile(gitconfigFile)
}

func (f *FileDbRepository) EditRepo(repoName string) error {
	// Check if repository exists
	_, err := f.GetRepositoryByName(repoName)
	if err != nil {
		return fmt.Errorf("repository '%s' not found: %w", repoName, err)
	}

	// Get repository db.[ext] file path
	repoDbFilePath, err := createRepoDbFilePath(f, repoName)
	if err != nil {
		return fmt.Errorf("failed to get repository path: %w", err)
	}

	return f.editFile(repoDbFilePath)
}

func (f *FileDbRepository) PushRepository(repoName string) error {
	// Check if repository exists
	_, err := f.GetRepositoryByName(repoName)
	if err != nil {
		return fmt.Errorf("repository '%s' not found: %w", repoName, err)
	}

	// Get repository path
	repoPath, err := f.DirectoryService.GetRepositoryPath(repoName)
	if err != nil {
		return fmt.Errorf("failed to get repository path: %w", err)
	}

	// Check if repository has git remote
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			return fmt.Errorf("repository '%s' is not a git repository", repoName)
		}
		return fmt.Errorf("failed to open git repository: %w", err)
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return fmt.Errorf("failed to get git remotes: %w", err)
	}

	if len(remotes) == 0 {
		return fmt.Errorf("repository '%s' does not have a git remote configured", repoName)
	}

	// Commit and push changes
	err = gitt.CommitAndPushChanges(repoPath)
	if err != nil {
		return fmt.Errorf("failed to push repository '%s': %w", repoName, err)
	}

	return nil
}

func (f *FileDbRepository) GetBasePath() (string, error) {
	basePath, err := f.getBasePath()
	if err != nil {
		return "", err
	}
	return basePath, nil
}

func (f *FileDbRepository) BonusInjection(enabledRepos []entity.Repository) (string, error) {
	for _, repo := range enabledRepos {
		if repo.GitConfigIncludePath == "" {
			continue
		}
		gitConfigPath, err := f.gitConfigPathProvider.GetPath()
		if err != nil {
			return "", err
		}
		err = gitconfig.AddNewIncludeIfNotExists(repo.GitConfigIncludePath, gitConfigPath)
		if err != nil {
			return "", err
		}
	}
	return "", nil
}

// ListRepoPath returns the base directory used for file-backed repositories
// and the full paths of its immediate subdirectories. The returned slice
// always includes the base path as the first element, followed by one
// entry for each direct child directory under that base path. An error is
// returned if the base path cannot be resolved or read.
func (f *FileDbRepository) ListRepoPath() ([]string, error) {
	path, err := f.getBasePath()
	if err != nil {
		return nil, err
	}
	repoPath := filepath.Join(path, "repositories")
	files, err := os.ReadDir(repoPath)
	if err != nil {
		return nil, err
	}
	paths := []string{}
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, filepath.Join(repoPath, file.Name()))
		}
	}
	return paths, nil
}

////////////////////
// Helper (internal) functions
////////////////////

func createRepoDbFilePath(f *FileDbRepository, name string) (string, error) {
	basePath, err := f.getBasePath()
	if err != nil {
		return "", err
	}
	fileName := "db." + f.fileHandler.Extension()
	repoDbFilePath := filepath.Join(basePath, "repositories", name, fileName)
	return repoDbFilePath, nil
}

func (f *FileDbRepository) getRepositoryGitconfigPath(name string) string {
	basePath, err := f.getBasePath()
	if err != nil {
		return ""
	}
	path := filepath.Join(basePath, "repositories", name, "gitconfig")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ""
	}
	return path
}

func (f *FileDbRepository) GetRepositoryByName(name string) (*entity.Repository, error) {
	repoPath, err := createRepoDbFilePath(f, name)
	if err != nil {
		return nil, err
	}
	repoDto, err := f.fileHandler.LoadRepositoryFile(repoPath)
	if err != nil {
		return nil, err
	}
	aliases := map[string]string{}
	if repoDto.Aliases == nil {
		aliases = map[string]string{}
	} else {
		aliases = repoDto.Aliases
	}
	exports := map[string]string{}
	if repoDto.Exports == nil {
		exports = map[string]string{}
	} else {
		exports = repoDto.Exports
	}

	repo := entity.Repository{
		Name:                 name,
		Aliases:              aliases,
		Exports:              exports,
		GitConfigIncludePath: f.getRepositoryGitconfigPath(name),
	}
	return &repo, nil
}

func (f *FileDbRepository) getUserPrefPath() (string, error) {
	basePath, err := f.getBasePath()
	if err != nil {
		return "", err
	}
	fileName := "user_preferences." + f.fileHandler.Extension()
	return filepath.Join(basePath, fileName), nil
}

func (f *FileDbRepository) getUserPreferences() (*UserPreferenceDto, error) {
	userPrefPath, err := f.getUserPrefPath()
	if err != nil {
		return nil, err
	}
	return f.fileHandler.LoadUserPreferenceFile(userPrefPath)
}

func (f *FileDbRepository) getBasePath() (string, error) {
	return f.PathProvider.GetPath()
}

func (f *FileDbRepository) GetAllRepositories() ([]entity.Repository, error) {
	allRepoNames, err := f.DirectoryService.ListRepositoryNames()
	if err != nil {
		return nil, err
	}
	allRepos := []entity.Repository{}
	for _, repoName := range allRepoNames {
		repo, err := f.GetRepositoryByName(repoName)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, *repo)
	}
	return allRepos, nil
}
