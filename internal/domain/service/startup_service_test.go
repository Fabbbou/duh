package service_test

import (
	"duh/internal/domain/repository"
	"duh/internal/domain/service"
	"duh/internal/domain/utils"
	"path/filepath"
	"testing"
)

func TestStartupService_Run(t *testing.T) {
	tempPath := t.TempDir()
	pathProvider := service.NewCustomPathProvider(tempPath)
	inMemrepoFactory := repository.NewMockDbRepositoryFactory()
	svc := service.NewStartupService(pathProvider, inMemrepoFactory)

	err := svc.Run()
	if err != nil {
		t.Errorf("StartupService.Run() error = %v, wantErr %v", err, false)
	}

	if utils.DirectoryExists(filepath.Join(tempPath, "repositories", "local")) == false {
		t.Errorf("Expected local repository directory to be created")
	}

	if utils.FileExists(filepath.Join(tempPath, "repositories", "local", "db.toml")) == false {
		t.Errorf("Expected local db.toml file to be created")
	}

	if utils.FileExists(filepath.Join(tempPath, "user_preferences.toml")) == false {
		t.Errorf("Expected user_preferences.toml file to be created")
	}

	if len(inMemrepoFactory.InmemoryCreated) == 0 {
		t.Errorf("Expected at least one in-memory repository to be created")
	}

	if len(inMemrepoFactory.InmemoryCreated) > 1 {
		t.Errorf("Expected only one in-memory repository to be created, got %d", len(inMemrepoFactory.InmemoryCreated))
	}

	repo := inMemrepoFactory.InmemoryCreated[0]
	values, err := repo.List(service.RepositoriesGroup)
	if err != nil {
		t.Errorf("Error listing values from user preference repository: %v", err)
	}

	if len(values) != 3 {
		t.Errorf("Expected 2 default user preference values, got %d", len(values))
	}
}
