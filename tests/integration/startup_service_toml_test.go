package integration_test

import (
	"duh/internal/domain/utils"
	"duh/internal/infrastructure/file_db"
	"path/filepath"
	"testing"
)

func TestInitDbService_Check(t *testing.T) {
	pathProvider := file_db.NewCustomPathProvider(t.TempDir())
	InitDbService := file_db.NewInitDbService(pathProvider)
	hasChanged, err := InitDbService.Check()
	if !hasChanged {
		t.Errorf("InitDbService.Check() expected to have changes on first run")
	}
	if err != nil {
		t.Fatalf("InitDbService.Check() returned an error: %v", err)
	}

	theDuhPath, err := pathProvider.GetPath()
	if err != nil {
		t.Fatalf("PathProvider.GetPath() returned an error: %v", err)
	}

	// Check if the duh directory was created
	if !utils.DirectoryExists(theDuhPath) {
		t.Errorf("Duh directory was not created at %s", theDuhPath)
	}
	// Check if the local repository directory was created
	localRepoPath := filepath.Join(theDuhPath, "repositories", "local")
	if !utils.DirectoryExists(localRepoPath) {
		t.Errorf("Local repository directory was not created at %s", localRepoPath)
	}
	// Check if the local repository db.toml file was created
	localDbPath := filepath.Join(localRepoPath, "db.toml")
	if !utils.FileExists(localDbPath) {
		t.Errorf("Local repository db.toml file was not created at %s", localDbPath)
	}
	// Check if the user_preferences.toml file was created
	userPrefPath := filepath.Join(theDuhPath, "user_preferences.toml")
	if !utils.FileExists(userPrefPath) {
		t.Errorf("User preferences file was not created at %s", userPrefPath)
	}
}
