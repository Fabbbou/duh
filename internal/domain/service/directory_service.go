package service

import (
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

func (ds *DirectoryService) CreateRepository(repositoryName string) (string, error) {
	basePath, err := ds.basePathProvider.GetPath()
	if err != nil {
		return "", err
	}
	repoPath := filepath.Join(basePath, "repositories", repositoryName)
	err = ds.ensureDirectoryExists(repoPath)
	if err != nil {
		return "", err
	}

	dbFilePath := filepath.Join(repoPath, "db.toml")
	_, err = os.Create(dbFilePath)
	if err != nil {
		return "", err
	}
	return repoPath, nil
}

// ListRepositoryNames returns a list of paths to existing repositories
// A repository is considered existing if it contains a "db.toml" file
func (ds *DirectoryService) ListRepositoryNames() ([]string, error) {
	basePath, err := ds.basePathProvider.GetPath()
	if err != nil {
		return nil, err
	}
	repoBasePath := filepath.Join(basePath, "repositories")
	entries, err := os.ReadDir(repoBasePath)
	if err != nil {
		return nil, err
	}

	var repoNames []string
	for _, entry := range entries {
		if entry.IsDir() {
			currentRepoPath := filepath.Join(repoBasePath, entry.Name())
			dbFile := filepath.Join(currentRepoPath, "db.toml")
			if os.Stat(dbFile); os.IsNotExist(err) {
				continue
			}
			repoNames = append(repoNames, entry.Name())
		}
	}
	return repoNames, nil
}
