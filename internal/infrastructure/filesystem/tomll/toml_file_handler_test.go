package tomll

import (
	"duh/internal/infrastructure/filesystem/common"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestTomlFileHandler_Extension(t *testing.T) {
	handler := &TomlFileHandler{}
	expected := "toml"
	actual := handler.Extension()
	if actual != expected {
		t.Errorf("Expected %s, but got %s", expected, actual)
	}
}

func TestTomlFileHandler_LoadRepositoryFile(t *testing.T) {
	handler := &TomlFileHandler{}

	// Create a temporary repository file
	tempDir := t.TempDir()
	repoFile := filepath.Join(tempDir, "test_repo.toml")

	repoContent := `[aliases]
build = "go build"
test = "go test"

[exports]
API_KEY = "test_key"
DEBUG = "true"

[metadata]
url_origin = "https://github.com/test/repo.git"
name_origin = "test-repo"`

	err := os.WriteFile(repoFile, []byte(repoContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	dto, err := handler.LoadRepositoryFile(repoFile)
	if err != nil {
		t.Fatalf("LoadRepositoryFile failed: %v", err)
	}

	// Verify aliases
	expectedAliases := map[string]string{
		"build": "go build",
		"test":  "go test",
	}
	if !reflect.DeepEqual(dto.Aliases, expectedAliases) {
		t.Errorf("Expected aliases %v, but got %v", expectedAliases, dto.Aliases)
	}

	// Verify exports
	expectedExports := map[string]string{
		"API_KEY": "test_key",
		"DEBUG":   "true",
	}
	if !reflect.DeepEqual(dto.Exports, expectedExports) {
		t.Errorf("Expected exports %v, but got %v", expectedExports, dto.Exports)
	}

	// Verify metadata
	if dto.Metadata.UrlOrigin != "https://github.com/test/repo.git" {
		t.Errorf("Expected UrlOrigin 'https://github.com/test/repo.git', but got '%s'", dto.Metadata.UrlOrigin)
	}
	if dto.Metadata.NameOrigin != "test-repo" {
		t.Errorf("Expected NameOrigin 'test-repo', but got '%s'", dto.Metadata.NameOrigin)
	}
}

func TestTomlFileHandler_SaveRepositoryFile(t *testing.T) {
	handler := &TomlFileHandler{}
	tempDir := t.TempDir()
	repoFile := filepath.Join(tempDir, "save_test_repo.toml")

	dto := &common.RepositoryDto{
		Aliases: map[string]string{
			"deploy": "docker deploy",
			"clean":  "make clean",
		},
		Exports: map[string]string{
			"ENV":     "production",
			"VERSION": "1.0.0",
		},
		Metadata: common.MetadataDto{
			UrlOrigin:  "https://github.com/example/project.git",
			NameOrigin: "example-project",
		},
	}

	err := handler.SaveRepositoryFile(repoFile, dto)
	if err != nil {
		t.Fatalf("SaveRepositoryFile failed: %v", err)
	}

	// Load it back to verify
	loadedDto, err := handler.LoadRepositoryFile(repoFile)
	if err != nil {
		t.Fatalf("Failed to load saved file: %v", err)
	}

	if !reflect.DeepEqual(dto.Aliases, loadedDto.Aliases) {
		t.Errorf("Saved and loaded aliases don't match")
	}
	if !reflect.DeepEqual(dto.Exports, loadedDto.Exports) {
		t.Errorf("Saved and loaded exports don't match")
	}
	if dto.Metadata != loadedDto.Metadata {
		t.Errorf("Saved and loaded metadata don't match")
	}
}

func TestTomlFileHandler_LoadUserPreferenceFile(t *testing.T) {
	handler := &TomlFileHandler{}
	tempDir := t.TempDir()
	userPrefFile := filepath.Join(tempDir, "user_pref.toml")

	userPrefContent := `[repositories]
activated_repos = ["repo1", "repo2", "repo3"]
default_repo_name = "repo1"`

	err := os.WriteFile(userPrefFile, []byte(userPrefContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	dto, err := handler.LoadUserPreferenceFile(userPrefFile)
	if err != nil {
		t.Fatalf("LoadUserPreferenceFile failed: %v", err)
	}

	expectedRepos := []string{"repo1", "repo2", "repo3"}
	if !reflect.DeepEqual(dto.Repositories.ActivatedRepositories, expectedRepos) {
		t.Errorf("Expected repositories %v, but got %v", expectedRepos, dto.Repositories.ActivatedRepositories)
	}

	if dto.Repositories.DefaultRepositoryName != "repo1" {
		t.Errorf("Expected default repository 'repo1', but got '%s'", dto.Repositories.DefaultRepositoryName)
	}
}

func TestTomlFileHandler_MigrateOldUserPrefFormat(t *testing.T) {
	handler := &TomlFileHandler{}

	// Use the existing old format file
	oldFormatFile := "old_user_pref_format.toml"

	dto, err := handler.LoadUserPreferenceFile(oldFormatFile)
	if err != nil {
		t.Fatalf("Failed to migrate old user preference file: %v", err)
	}

	// Check that migration worked correctly
	expectedRepos := []string{"repo1", "repo2"}
	if !reflect.DeepEqual(dto.Repositories.ActivatedRepositories, expectedRepos) {
		t.Errorf("Expected migrated repositories %v, but got %v", expectedRepos, dto.Repositories.ActivatedRepositories)
	}

	if dto.Repositories.DefaultRepositoryName != "repo1" {
		t.Errorf("Expected migrated default repository 'repo1', but got '%s'", dto.Repositories.DefaultRepositoryName)
	}
}

func TestTomlFileHandler_SaveUserPreferenceFile(t *testing.T) {
	handler := &TomlFileHandler{}
	tempDir := t.TempDir()
	userPrefFile := filepath.Join(tempDir, "save_user_pref.toml")

	dto := &common.UserPreferenceDto{
		Repositories: common.RepositoriesPreferenceDto{
			ActivatedRepositories: []string{"main", "dev", "staging"},
			DefaultRepositoryName: "main",
		},
	}

	err := handler.SaveUserPreferenceFile(userPrefFile, dto)
	if err != nil {
		t.Fatalf("SaveUserPreferenceFile failed: %v", err)
	}

	// Load it back to verify
	loadedDto, err := handler.LoadUserPreferenceFile(userPrefFile)
	if err != nil {
		t.Fatalf("Failed to load saved user preference file: %v", err)
	}

	if !reflect.DeepEqual(dto.Repositories.ActivatedRepositories, loadedDto.Repositories.ActivatedRepositories) {
		t.Errorf("Saved and loaded activated repositories don't match")
	}
	if dto.Repositories.DefaultRepositoryName != loadedDto.Repositories.DefaultRepositoryName {
		t.Errorf("Saved and loaded default repository name don't match")
	}
}

func TestTomlFileHandler_LoadRepositoryFile_FileNotFound(t *testing.T) {
	handler := &TomlFileHandler{}

	_, err := handler.LoadRepositoryFile("nonexistent_file.toml")
	if err == nil {
		t.Error("Expected error for nonexistent file, but got none")
	}
}

func TestTomlFileHandler_LoadUserPreferenceFile_FileNotFound(t *testing.T) {
	handler := &TomlFileHandler{}

	_, err := handler.LoadUserPreferenceFile("nonexistent_user_pref.toml")
	if err == nil {
		t.Error("Expected error for nonexistent file, but got none")
	}
}
