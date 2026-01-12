package common

import (
	"duh/internal/domain/constants"
	"os"
	"path/filepath"
)

type DirectoryService struct {
	basePathProvider PathProvider
}

func NewDirectoryService(basePathProvider PathProvider) *DirectoryService {
	return &DirectoryService{
		basePathProvider: basePathProvider,
	}
}

func (ds *DirectoryService) ensureDirectoryExists(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func (ds *DirectoryService) CreateGitconfigFile(repositoryName string) error {
	repoPath, err := ds.GetRepositoryPath(repositoryName)
	if err != nil {
		return err
	}

	// create a gitconfig file for the repository
	gitconfigPath := filepath.Join(repoPath, "gitconfig")
	gitconfigFile, err := os.Create(gitconfigPath)
	if err != nil {
		return err
	}
	_, err = gitconfigFile.WriteString("[alias]\n\t")
	if err != nil {
		return err
	}
	return gitconfigFile.Close()
}

func (ds *DirectoryService) CreatePackage(repositoryName string) (string, error) {
	repoPath, err := ds.GetRepositoryPath(repositoryName)
	if err != nil {
		return "", err
	}

	err = ds.CreateGitconfigFile(repositoryName)
	if err != nil {
		return "", err
	}

	dbFilePath := filepath.Join(repoPath, constants.PackageDbFileName+".toml")
	if _, err := os.Stat(dbFilePath); os.IsExist(err) {
		return repoPath, nil
	}
	file, err := os.Create(dbFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return repoPath, nil
}

func (ds *DirectoryService) DeletePackage(repositoryName string) error {
	repoPath, err := ds.GetRepositoryPath(repositoryName)
	if err != nil {
		return err
	}
	err = os.RemoveAll(repoPath)
	if err != nil {
		return err
	}
	return nil
}

// ListRepositoryNames returns a list of paths to existing repositories
// A repository is considered existing if it contains a "db.toml" file
func (ds *DirectoryService) ListRepositoryNames() ([]string, error) {
	basePath, err := ds.basePathProvider.GetPath()
	if err != nil {
		return nil, err
	}
	repoBasePath := filepath.Join(basePath, constants.PackagesDirName)
	entries, err := os.ReadDir(repoBasePath)
	if err != nil {
		return nil, err
	}

	var repoNames []string
	for _, entry := range entries {
		if entry.IsDir() {
			currentRepoPath := filepath.Join(repoBasePath, entry.Name())
			dbFile := filepath.Join(currentRepoPath, constants.PackageDbFileName+".toml")
			if os.Stat(dbFile); os.IsNotExist(err) {
				continue
			}
			repoNames = append(repoNames, entry.Name())
		}
	}
	return repoNames, nil
}

func (ds *DirectoryService) GetRepositoryPath(repositoryName string) (string, error) {
	basePath, err := ds.basePathProvider.GetPath()
	if err != nil {
		return "", err
	}
	repoPath := filepath.Join(basePath, constants.PackagesDirName, repositoryName)
	err = ds.ensureDirectoryExists(repoPath)
	if err != nil {
		return "", err
	}
	return repoPath, nil
}
