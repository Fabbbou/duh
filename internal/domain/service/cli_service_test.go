package service

import (
	"duh/internal/domain/entity"
	"duh/internal/domain/repository"
	"strings"
	"testing"
)

var expectedInjection = []string{
	`alias ll="ls -la"`,
	`alias gs="git status"`,
	`alias ca="echo \"Complex Alias\""`,
	`export PATH="/usr/local/bin:$PATH"`,
	`export GOENV="development"`,
	`alias 2ll="2ls -la"`,
	`alias 2gs="2git status"`,
	`alias 2ca="2echo \"Complex Alias\""`,
	`export 2PATH="2/usr/local/bin:$PATH"`,
	`export 2GOENV="2development"`,
}

func TestCliService_Inject(t *testing.T) {
	mockRepo1 := repository.NewMockInmemoryDbRepository()
	mockRepo1.Upsert(entity.Aliases, "ll", "ls -la")
	mockRepo1.Upsert(entity.Aliases, "gs", "git status")
	mockRepo1.Upsert(entity.Aliases, "ca", `echo "Complex Alias"`)
	mockRepo1.Upsert(entity.Exports, "PATH", "/usr/local/bin:$PATH")
	mockRepo1.Upsert(entity.Exports, "GOENV", "development")

	mockRepo2 := repository.NewMockInmemoryDbRepository()
	mockRepo2.Upsert(entity.Aliases, "2ll", "2ls -la")
	mockRepo2.Upsert(entity.Aliases, "2gs", "2git status")
	mockRepo2.Upsert(entity.Aliases, "2ca", `2echo "Complex Alias"`)
	mockRepo2.Upsert(entity.Exports, "2PATH", "2/usr/local/bin:$PATH")
	mockRepo2.Upsert(entity.Exports, "2GOENV", "2development")

	mockRepo3 := repository.NewMockInmemoryDbRepository()
	mockRepo3.Upsert(entity.Aliases, "3ll", "3ls -la")
	mockRepo3.Upsert(entity.Aliases, "3gs", "3git status")
	mockRepo3.Upsert(entity.Aliases, "3ca", `3echo "Complex Alias"`)
	mockRepo3.Upsert(entity.Exports, "3PATH", "3/usr/local/bin:$PATH")
	mockRepo3.Upsert(entity.Exports, "3GOENV", "3development")

	repositoryService := NewRepositoriesService(nil, nil)
	repositoryService.allRepositories["repo1"] = mockRepo1
	repositoryService.allRepositories["repo2"] = mockRepo2
	repositoryService.allRepositories["repo3"] = mockRepo3

	mockRepo4 := repository.NewMockInmemoryDbRepository()
	mockRepo4.Upsert(RepositoriesGroup, ActivatedRepositoriesKey, "repo1,repo2")
	userPreferenceService := NewUserPreferenceService(mockRepo4)

	cliService := NewCliService(repositoryService, userPreferenceService)

	injection, err := cliService.Inject()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	for _, expectedLine := range expectedInjection {
		if !containsLine(injection, expectedLine) {
			t.Errorf("Expected injection to contain line: %s", expectedLine)
		}
	}
}

func containsLine(injection string, line string) bool {
	lines := strings.Split(injection, "\n")
	for _, l := range lines {
		if l == line {
			return true
		}
	}
	return false
}

func TestCliService_AddAndRemoveAlias(t *testing.T) {
	mockRepoDb := repository.NewMockInmemoryDbRepository()
	mockRepoUser := repository.NewMockInmemoryDbRepository()

	repositoryService := NewRepositoriesService(nil, nil)
	repositoryService.allRepositories["local"] = mockRepoDb

	userPreferenceService := NewUserPreferenceService(mockRepoUser)
	err := userPreferenceService.InitUserPreference()
	if err != nil {
		t.Fatalf("Expected no error on InitUserPreference, got %v", err)
	}

	cliService := NewCliService(repositoryService, userPreferenceService)

	err = cliService.SetAlias("testalias", "testcommand")
	if err != nil {
		t.Fatalf("Expected no error on SetAlias, got %v", err)
	}

	aliases, err := cliService.ListAliases()
	if err != nil {
		t.Fatalf("Expected no error on ListAliases, got %v", err)
	}

	if val, exists := aliases["testalias"]; !exists || val != "testcommand" {
		t.Errorf("Expected alias 'testalias' with value 'testcommand', got %v", aliases)
	}

	err = cliService.RemoveAlias("testalias")
	if err != nil {
		t.Fatalf("Expected no error on RemoveAlias, got %v", err)
	}

	aliases, err = cliService.ListAliases()
	if err != nil {
		t.Fatalf("Expected no error on ListAliases after removal, got %v", err)
	}

	if _, exists := aliases["testalias"]; exists {
		t.Errorf("Expected alias 'testalias' to be removed, but it still exists")
	}
}
