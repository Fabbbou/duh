package integration

import (
	"duh/internal/domain/constants"
	"duh/internal/infrastructure/filesystem/common"
	"duh/internal/infrastructure/filesystem/fs_function_adapter"
	"duh/internal/infrastructure/filesystem/fs_user_repository"
	"duh/internal/infrastructure/filesystem/tomll"
	"os"
	"path/filepath"
	"testing"
)

func TestFSFunctionsRepository_GetActivatedScripts_Integration(t *testing.T) {
	// Create temporary directory structure
	tempDir := t.TempDir()
	reposDir := filepath.Join(tempDir, constants.PackagesDirName)

	// Create repository structures
	repo1Dir := filepath.Join(reposDir, "repo1", constants.PackageFunctionsDirName)
	repo2Dir := filepath.Join(reposDir, "repo2", constants.PackageFunctionsDirName)

	err := os.MkdirAll(repo1Dir, 0755)
	if err != nil {
		t.Fatalf("Failed to create repo1 functions dir: %v", err)
	}

	err = os.MkdirAll(repo2Dir, 0755)
	if err != nil {
		t.Fatalf("Failed to create repo2 functions dir: %v", err)
	}

	// Create user preferences directory and file
	userPrefPath := filepath.Join(tempDir, constants.DuhConfigFileName+"."+"toml")
	userPrefContent := `[repositories]
activated_repos = ["repo1", "repo2"]
default_repo_name = "repo1"`

	err = os.WriteFile(userPrefPath, []byte(userPrefContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create user preferences file: %v", err)
	}

	// Create script files in repo1
	repo1Script1 := `#!/bin/bash

# Deploy application to production
# Parameters: version (string)
# Returns: 0 on success, 1 on failure
function deploy() {
    local version="$1"
    echo "Deploying version $version"
    return 0
}

# Check application health
function health_check() {
    echo "Application is healthy"
    return 0
}`

	script1Path := filepath.Join(repo1Dir, "deploy.sh")
	err = os.WriteFile(script1Path, []byte(repo1Script1), 0644)
	if err != nil {
		t.Fatalf("Failed to create script1: %v", err)
	}

	repo1Script2 := `#!/bin/bash

# Backup database
# Parameters: backup_name (string)
function backup_db() {
    local name="$1"
    echo "Backing up database to $name"
}`

	script2Path := filepath.Join(repo1Dir, "backup.sh")
	err = os.WriteFile(script2Path, []byte(repo1Script2), 0644)
	if err != nil {
		t.Fatalf("Failed to create script2: %v", err)
	}

	// Create script files in repo2
	repo2Script := `#!/bin/bash

# Monitor system resources
# Parameters: none
# Returns: resource usage info
function monitor() {
    echo "System monitoring"
    top -n 1
}

# Some global code that shouldn't be here
echo "Global execution in repo2"`

	script3Path := filepath.Join(repo2Dir, "monitor.sh")
	err = os.WriteFile(script3Path, []byte(repo2Script), 0644)
	if err != nil {
		t.Fatalf("Failed to create script3: %v", err)
	}

	// Set up dependencies - use actual implementations
	pathProvider := common.NewCustomPathProvider(tempDir)
	fileHandler := &tomll.TomlFileHandler{}
	userPrefRepo := fs_user_repository.NewFsUserRepository(fileHandler, pathProvider)

	// Create repository under test
	repo := fs_function_adapter.NewFSFunctionsRepository(pathProvider, userPrefRepo)

	// Execute test
	scripts, err := repo.GetActivatedScripts()
	if err != nil {
		t.Fatalf("GetActivatedScripts failed: %v", err)
	}

	// Verify results
	if len(scripts) != 3 {
		t.Errorf("Expected 3 scripts, got %d", len(scripts))
	}

	// Verify script names and basic properties
	scriptNames := make(map[string]bool)
	for _, script := range scripts {
		scriptNames[script.Name] = true

		// Verify basic properties
		if script.PathToFile == "" {
			t.Errorf("Script %s should have PathToFile", script.Name)
		}
		if script.DataToInject == "" {
			t.Errorf("Script %s should have DataToInject", script.Name)
		}

		// Verify functions are detected
		switch script.Name {
		case "deploy":
			if len(script.Functions) != 2 {
				t.Errorf("Deploy script should have 2 functions, got %d", len(script.Functions))
			}
		case "backup":
			if len(script.Functions) != 1 {
				t.Errorf("Backup script should have 1 function, got %d", len(script.Functions))
			}
		case "monitor":
			if len(script.Functions) != 1 {
				t.Errorf("Monitor script should have 1 function, got %d", len(script.Functions))
			}
			if len(script.Warnings) == 0 {
				t.Error("Monitor script should have warnings due to global code")
			}
		}
	}

	// Verify all expected scripts are present
	expectedScripts := []string{"deploy", "backup", "monitor"}
	for _, expected := range expectedScripts {
		if !scriptNames[expected] {
			t.Errorf("Expected to find script %s", expected)
		}
	}
}

func TestFSFunctionsRepository_GetActivatedScripts_EmptyActivatedRepos(t *testing.T) {
	tempDir := t.TempDir()

	// Create user preferences with no activated repos
	userPrefPath := filepath.Join(tempDir, constants.DuhConfigFileName+"."+"toml")
	userPrefContent := `[repositories]
activated_repos = []
default_repo_name = ""`

	err := os.WriteFile(userPrefPath, []byte(userPrefContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create user preferences file: %v", err)
	}

	pathProvider := common.NewCustomPathProvider(tempDir)
	fileHandler := &tomll.TomlFileHandler{}
	userPrefRepo := fs_user_repository.NewFsUserRepository(fileHandler, pathProvider)

	repo := fs_function_adapter.NewFSFunctionsRepository(pathProvider, userPrefRepo)

	scripts, err := repo.GetActivatedScripts()
	if err != nil {
		t.Fatalf("GetActivatedScripts should not error with empty repos: %v", err)
	}

	if len(scripts) != 0 {
		t.Errorf("Expected 0 scripts with no activated repos, got %d", len(scripts))
	}
}

func TestFSFunctionsRepository_GetActivatedScripts_MissingRepositoryDirectories(t *testing.T) {
	tempDir := t.TempDir()

	// Create user preferences with non-existent repos
	userPrefPath := filepath.Join(tempDir, constants.DuhConfigFileName+"."+"toml")
	userPrefContent := `[repositories]
activated_repos = ["nonexistent1", "nonexistent2"]
default_repo_name = "nonexistent1"`

	err := os.WriteFile(userPrefPath, []byte(userPrefContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create user preferences file: %v", err)
	}

	pathProvider := common.NewCustomPathProvider(tempDir)
	fileHandler := &tomll.TomlFileHandler{}
	userPrefRepo := fs_user_repository.NewFsUserRepository(fileHandler, pathProvider)

	repo := fs_function_adapter.NewFSFunctionsRepository(pathProvider, userPrefRepo)

	scripts, err := repo.GetActivatedScripts()
	if err != nil {
		t.Fatalf("GetActivatedScripts should not error with missing directories: %v", err)
	}

	// Should return empty list when directories don't exist
	if len(scripts) != 0 {
		t.Errorf("Expected 0 scripts with missing directories, got %d", len(scripts))
	}
}
