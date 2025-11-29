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

func (ds *DirectoryService) CreateRepositoryPath(repositoryName string) (string, error) {
	basePath, err := ds.basePathProvider.GetPath()
	if err != nil {
		return "", err
	}
	repoPath := filepath.Join(basePath, "repositories", repositoryName)
	err = ds.ensureDirectoryExists(repoPath)
	if err != nil {
		return "", err
	}
	return repoPath, nil
}
