package integration

import (
	"duh/internal/domain/entity"
	"duh/internal/domain/utils"
	"duh/internal/infrastructure/filesystem/common"
	"duh/internal/infrastructure/filesystem/file_db"
	"duh/internal/infrastructure/filesystem/tomll"
	"os"
	"path/filepath"
	"testing"
)

func TestInitDbService_Check(t *testing.T) {
	tempPath := t.TempDir()
	defer os.RemoveAll(tempPath)
	pathProvider := common.NewCustomPathProvider(tempPath)
	svc := file_db.NewInitDbService(pathProvider, &tomll.TomlFileHandler{})
	hasChanged, err := svc.Check()
	if !hasChanged {
		t.Errorf("InitDbService.Check() expected to have changes on first run")
	}
	if err != nil {
		t.Errorf("InitDbService.Check() error = %v, wantErr %v", err, false)
	}
	if err != nil {
		t.Errorf("InitDbService.Check() error = %v, wantErr %v", err, false)
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

	userPrefs, err := tomll.LoadToml[tomll.UserPreferenceToml](filepath.Join(tempPath, "user_preferences.toml"))
	if err != nil {
		t.Errorf("Error retrieving user preferences: %v", err)
	}
	if userPrefs.Repositories.ActivatedRepositories == nil {
		t.Errorf("Expected ActivatedRepositories to be initialized")
	}
	expectedRepos := []entity.Repository{{
		Name: "local",
	}}
	if len(userPrefs.Repositories.ActivatedRepositories) != len(expectedRepos) {
		t.Errorf("Expected %d activated repositories, got %d", len(expectedRepos), len(userPrefs.Repositories.ActivatedRepositories))
	} else {
		for i, repo := range userPrefs.Repositories.ActivatedRepositories {
			if repo != expectedRepos[i].Name {
				t.Errorf("Expected repository name %s, got %s", expectedRepos[i].Name, repo)
			}
		}
	}
}
