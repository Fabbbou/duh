package file_service

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
	repoPath, err := directoryService.CreateRepository(repoName)
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
	dbFilePath := filepath.Join(repoPath, "db.toml")
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		t.Fatalf("Expected db.toml file to exist at %s", dbFilePath)
	}
}

func TestCreateRepositoryDirectory_Error(t *testing.T) {
	invalidPathProvider := NewCustomPathProvider("/invalid\000path")
	directoryService := NewDirectoryService(invalidPathProvider)

	_, err := directoryService.CreateRepository("test-repo")
	if err == nil {
		t.Fatalf("Expected an error due to invalid path, got nil")
	}
}

func TestListRepositoryNames(t *testing.T) {
	tempDir := t.TempDir()
	customPathProvider := NewCustomPathProvider(tempDir)
	directoryService := NewDirectoryService(customPathProvider)

	expectedRepoNames := []string{"repo1", "repo2", "repo3"}
	for _, name := range expectedRepoNames {
		_, err := directoryService.CreateRepository(name)
		if err != nil {
			t.Fatalf("Expected no error creating repository %s, got %v", name, err)
		}
	}

	createdRepoNames, err := directoryService.ListRepositoryNames()
	if err != nil {
		t.Fatalf("Expected no error listing repository paths, got %v", err)
	}

	if len(createdRepoNames) != len(expectedRepoNames) {
		t.Fatalf("Expected %d repository paths, got %d", len(expectedRepoNames), len(createdRepoNames))
	}

	for _, name := range expectedRepoNames {
		found := false
		for _, path := range createdRepoNames {
			if path == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find repository path %s", name)
		}
	}
}
