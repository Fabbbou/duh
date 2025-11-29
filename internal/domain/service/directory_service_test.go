package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateRepositoryDirectory(t *testing.T) {
	tempDir := t.TempDir()
	customPathProvider := NewCustomPathProvider(tempDir)
	directoryService := NewDirectoryService(customPathProvider)

	repoName := "test-repo"
	repoPath, err := directoryService.CreateRepositoryPath(repoName)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedPath := filepath.Join(tempDir, "repositories", repoName)
	if repoPath != expectedPath {
		t.Errorf("Expected path %s, got %s", expectedPath, repoPath)
	}

	info, err := os.Stat(repoPath)
	if os.IsNotExist(err) {
		t.Fatalf("Expected directory to exist at %s", repoPath)
	}
	if !info.IsDir() {
		t.Fatalf("Expected a directory at %s", repoPath)
	}
}

func TestCreateRepositoryDirectory_Error(t *testing.T) {
	invalidPathProvider := NewCustomPathProvider("/invalid\000path")
	directoryService := NewDirectoryService(invalidPathProvider)

	_, err := directoryService.CreateRepositoryPath("test-repo")
	if err == nil {
		t.Fatalf("Expected an error due to invalid path, got nil")
	}
}
