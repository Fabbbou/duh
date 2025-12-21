package file_db

import (
	"duh/internal/domain/entity"
	"duh/internal/domain/utils"
	"duh/internal/infrastructure/file_db/toml_repo"
	"os"
	"path/filepath"
	"testing"
)

func TestInitDbService_Run(t *testing.T) {
	tempPath := t.TempDir()
	defer os.RemoveAll(tempPath)
	pathProvider := NewCustomPathProvider(tempPath)
	svc := NewInitDbService(pathProvider)

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

	userPrefs, err := toml_repo.LoadUserPreferences(filepath.Join(tempPath, "user_preferences.toml"))
	if err != nil {
		t.Errorf("Error retrieving user preferences: %v", err)
	}
	if userPrefs.GetActivatedRepositories() == nil {
		t.Errorf("Expected ActivatedRepositories to be initialized")
	}
	expectedRepos := []entity.Repository{{
		Name: "local",
	}}
	if len(userPrefs.GetActivatedRepositories()) != len(expectedRepos) {
		t.Errorf("Expected %d activated repositories, got %d", len(expectedRepos), len(userPrefs.GetActivatedRepositories()))
	} else {
		for i, repo := range userPrefs.GetActivatedRepositories() {
			if repo != expectedRepos[i].Name {
				t.Errorf("Expected repository name %s, got %s", expectedRepos[i].Name, repo)
			}
		}
	}
}
