package service

import (
	"duh/internal/domain/repository"
	"testing"
)

func TestUserPreferenceService_GetCurrentActivatedRepositories(t *testing.T) {
	userDbRepo := repository.NewMockInmemoryDbRepository()
	userDbRepo.Upsert(RepositoriesGroup, ActivatedRepositoriesKey, "repo1,repo2,repo3")

	userPrefService := NewUserPreferenceService(userDbRepo)
	activatedRepos, err := userPrefService.GetCurrentActivatedRepositories()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedRepos := []string{"repo1", "repo2", "repo3"}
	if len(activatedRepos) != len(expectedRepos) {
		t.Fatalf("expected %d repositories, got %d", len(expectedRepos), len(activatedRepos))
	}
	for i, repo := range activatedRepos {
		if repo != expectedRepos[i] {
			t.Errorf("expected repo %s at index %d, got %s", expectedRepos[i], i, repo)
		}
	}
}

func TestUserPreferenceService_AddActivatedRepository(t *testing.T) {
	userDbRepo := repository.NewMockInmemoryDbRepository()
	userPrefService := NewUserPreferenceService(userDbRepo)

	err := userPrefService.AddActivatedRepository("repo1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = userPrefService.AddActivatedRepository("repo2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	activatedRepos, err := userPrefService.GetCurrentActivatedRepositories()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedRepos := []string{"repo1", "repo2"}
	if len(activatedRepos) != len(expectedRepos) {
		t.Fatalf("expected %d repositories, got %d", len(expectedRepos), len(activatedRepos))
	}
	for i, repo := range activatedRepos {
		if repo != expectedRepos[i] {
			t.Errorf("expected repo %s at index %d, got %s", expectedRepos[i], i, repo)
		}
	}

	// Test adding a duplicate repository
	err = userPrefService.AddActivatedRepository("repo1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	activatedRepos, err = userPrefService.GetCurrentActivatedRepositories()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(activatedRepos) != len(expectedRepos) {
		t.Fatalf("expected %d repositories after adding duplicate, got %d", len(expectedRepos), len(activatedRepos))
	}
}

func TestUserPreferenceService_RemoveActivatedRepository(t *testing.T) {
	userDbRepo := repository.NewMockInmemoryDbRepository()
	userDbRepo.Upsert(RepositoriesGroup, ActivatedRepositoriesKey, "repo1,repo2,repo3")

	userPrefService := NewUserPreferenceService(userDbRepo)

	err := userPrefService.RemoveActivatedRepository("repo2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	activatedRepos, err := userPrefService.GetCurrentActivatedRepositories()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedRepos := []string{"repo1", "repo3"}
	if len(activatedRepos) != len(expectedRepos) {
		t.Fatalf("expected %d repositories, got %d", len(expectedRepos), len(activatedRepos))
	}
	for i, repo := range activatedRepos {
		if repo != expectedRepos[i] {
			t.Errorf("expected repo %s at index %d, got %s", expectedRepos[i], i, repo)
		}
	}
}

func TestUserPreferenceService_SetAndGetLocalRepoName(t *testing.T) {
	userDbRepo := repository.NewMockInmemoryDbRepository()
	userPrefService := NewUserPreferenceService(userDbRepo)

	err := userPrefService.SetLocalRepoName("my-local-repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	localRepoName, err := userPrefService.GetLocalRepoName()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedName := "my-local-repo"
	if localRepoName != expectedName {
		t.Errorf("expected local repo name %s, got %s", expectedName, localRepoName)
	}
}

func TestUserPreferenceService_InitUserPreference(t *testing.T) {
	userDbRepo := repository.NewMockInmemoryDbRepository()
	userPrefService := NewUserPreferenceService(userDbRepo)

	err := userPrefService.InitUserPreference()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	localRepoName, err := userPrefService.GetLocalRepoName()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedName := "local"
	if localRepoName != expectedName {
		t.Errorf("expected local repo name %s, got %s", expectedName, localRepoName)
	}

	activatedRepos, err := userPrefService.GetCurrentActivatedRepositories()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	defaultRepoName, err := userPrefService.GetDefaultRepoName()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if defaultRepoName != "local" {
		t.Errorf("expected default repo name to be 'local', got %s", defaultRepoName)
	}

	expectedRepos := []string{"local"}
	if len(activatedRepos) != len(expectedRepos) {
		t.Fatalf("expected %d repositories, got %d", len(expectedRepos), len(activatedRepos))
	}
	for i, repo := range activatedRepos {
		if repo != expectedRepos[i] {
			t.Errorf("expected repo %s at index %d, got %s", expectedRepos[i], i, repo)
		}
	}

	// Test initializing again should return an error
	err = userPrefService.InitUserPreference()
	if err == nil {
		t.Fatalf("expected error when initializing user preference again, got nil")
	}
}
