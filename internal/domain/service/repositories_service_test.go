package service

import (
	"duh/internal/domain/repository"
	"testing"
)

func TestLoadRepositories(t *testing.T) {
	// Setup mock DirectoryService
	customDirService := &DirectoryService{
		basePathProvider: &CustomPathProvider{customPath: t.TempDir()},
	}
	customDirService.CreateRepository("repo1")
	customDirService.CreateRepository("repo2")

	// Setup mock DbRepositoryFactory
	mockDbRepoFactory := repository.NewMockDbRepositoryFactory()

	// Create RepositoriesService
	repoService := NewRepositoriesService(customDirService, mockDbRepoFactory)

	// Call LoadExistingRepositories
	err := repoService.LoadExistingRepositories()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify that repositories were loaded correctly
	expectedRepoCount := 2 // Assuming there are 2 mock repositories
	if len(repoService.allRepositories) != expectedRepoCount {
		t.Errorf("Expected %d repositories, got %d", expectedRepoCount, len(repoService.allRepositories))
	}

	_, exists := repoService.allRepositories["repo1"]
	if !exists {
		t.Errorf("Expected repository 'repo1' to be loaded")
	}

	_, exists = repoService.allRepositories["repo2"]
	if !exists {
		t.Errorf("Expected repository 'repo2' to be loaded")
	}
}
