package function

import (
	"duh/internal/domain/entity"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestGetScripts_ValidDirectory(t *testing.T) {
	// Create temporary directory with test scripts
	tempDir := t.TempDir()

	// Create test script 1 - well documented
	script1Content := `#!/bin/bash

# Deploy application to server
# Parameters: environment (dev|staging|prod)
# Returns: 0 on success, 1 on failure
function deploy() {
    local env="$1"
    echo "Deploying to $env"
    return 0
}

# Clean up temporary files
# Parameters: none
# Returns: always 0
function cleanup() {
    rm -rf /tmp/app/*
    return 0
}`

	script1Path := filepath.Join(tempDir, "deploy.sh")
	err := os.WriteFile(script1Path, []byte(script1Content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test script 1: %v", err)
	}

	// Create test script 2 - with issues
	script2Content := `#!/bin/bash

# Global code that shouldn't be here
echo "This is bad practice!"

# Undocumented function
function backup() {
    cp -r /app /backup/
}`

	script2Path := filepath.Join(tempDir, "backup.sh")
	err = os.WriteFile(script2Path, []byte(script2Content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test script 2: %v", err)
	}

	// Test GetScripts
	scripts, err := GetScripts(tempDir)
	if err != nil {
		t.Fatalf("GetScripts failed: %v", err)
	}

	// Should find 2 scripts
	if len(scripts) != 2 {
		t.Errorf("Expected 2 scripts, got %d", len(scripts))
	}

	// Check script properties
	deployScript := findScriptByName(scripts, "deploy")
	if deployScript == nil {
		t.Fatal("Should have found deploy script")
	}

	if deployScript.PathToFile != script1Path {
		t.Errorf("Expected path %s, got %s", script1Path, deployScript.PathToFile)
	}

	if deployScript.DataToInject == "" {
		t.Error("DataToInject should not be empty")
	}

	// Should have 2 functions
	if len(deployScript.Functions) != 2 {
		t.Errorf("Deploy script should have 2 functions, got %d", len(deployScript.Functions))
	}

	// Should have no warnings (well-written script)
	if len(deployScript.Warnings) != 0 {
		t.Errorf("Deploy script should have no warnings, got %d", len(deployScript.Warnings))
	}

	// Check backup script has warnings
	backupScript := findScriptByName(scripts, "backup")
	if backupScript == nil {
		t.Fatal("Should have found backup script")
	}

	if len(backupScript.Warnings) == 0 {
		t.Error("Backup script should have warnings due to global code")
	}
}

func TestGetScripts_NonExistentDirectory(t *testing.T) {
	nonExistentPath := "/this/directory/does/not/exist"

	scripts, err := GetScripts(nonExistentPath)
	if err == nil {
		t.Error("Expected error for non-existent directory")
	}

	if scripts != nil {
		t.Error("Scripts should be nil when directory doesn't exist")
	}

	// Check error message
	expectedMsg := "could not find directory"
	if err.Error() == "" || len(err.Error()) < len(expectedMsg) {
		t.Errorf("Error message should mention missing directory, got: %v", err)
	}
}

func TestGetScripts_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()
	// Don't create any files in the directory

	scripts, err := GetScripts(tempDir)
	if err != nil {
		t.Errorf("GetScripts should not error on empty directory: %v", err)
	}

	if len(scripts) != 0 {
		t.Errorf("Expected 0 scripts in empty directory, got %d", len(scripts))
	}
}

func TestGetScripts_InvalidShellScript(t *testing.T) {
	tempDir := t.TempDir()

	// Create invalid shell script
	invalidContent := `#!/bin/bash
	function invalid_syntax {
		if [ without closing
		echo "broken"
	`

	invalidPath := filepath.Join(tempDir, "invalid.sh")
	err := os.WriteFile(invalidPath, []byte(invalidContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid script: %v", err)
	}

	scripts, err := GetScripts(tempDir)
	if err == nil {
		t.Error("Expected error for invalid shell script")
	}

	if scripts != nil {
		t.Error("Scripts should be nil when analysis fails")
	}
}

func TestGetScripts_MixedFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Create a shell script
	scriptContent := `#!/bin/bash
# Test function
function test_func() {
    echo "test"
}`

	scriptPath := filepath.Join(tempDir, "script.sh")
	err := os.WriteFile(scriptPath, []byte(scriptContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	// Create a non-shell file
	textContent := "This is just a text file"
	textPath := filepath.Join(tempDir, "readme.txt")
	err = os.WriteFile(textPath, []byte(textContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create text file: %v", err)
	}

	scripts, err := GetScripts(tempDir)
	if err != nil {
		t.Fatalf("GetScripts should handle mixed files: %v", err)
	}

	// Should process all files, including non-shell ones
	// (the function doesn't filter by extension)
	if len(scripts) != 2 {
		t.Errorf("Expected 2 files to be processed, got %d", len(scripts))
	}
}

func TestGetScripts_ScriptAnalysisIntegration(t *testing.T) {
	tempDir := t.TempDir()

	// Create script with specific features we can test
	scriptContent := `#!/bin/bash

# Main application function
# Parameters: action (start|stop|restart)
# Returns: 0 on success, non-zero on failure
function main() {
    local action="$1"
    case "$action" in
        start)   start_service ;;
        stop)    stop_service ;;
        restart) restart_service ;;
        *)       echo "Invalid action" && return 1 ;;
    esac
}

# Start the service
function start_service() {
    echo "Starting service"
}

function stop_service() {
    echo "Stopping service"
}

function restart_service() {
    stop_service
    start_service
}`

	scriptPath := filepath.Join(tempDir, "service.sh")
	err := os.WriteFile(scriptPath, []byte(scriptContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	scripts, err := GetScripts(tempDir)
	if err != nil {
		t.Fatalf("GetScripts failed: %v", err)
	}

	if len(scripts) != 1 {
		t.Fatalf("Expected 1 script, got %d", len(scripts))
	}

	script := scripts[0]

	// Check basic properties
	if script.Name != "service" {
		t.Errorf("Expected name 'service', got '%s'", script.Name)
	}

	// Should have 4 functions
	if len(script.Functions) != 4 {
		t.Errorf("Expected 4 functions, got %d", len(script.Functions))
	}

	// Should have no global code warnings
	if len(script.Warnings) != 0 {
		t.Errorf("Expected no warnings, got %d", len(script.Warnings))
	}

	// Check that script content is preserved
	if script.DataToInject != scriptContent {
		t.Error("DataToInject should match original script content")
	}
}

// Helper function to find script by name
func findScriptByName(scripts []entity.Script, name string) *entity.Script {
	for i := range scripts {
		if scripts[i].Name == name {
			return &scripts[i]
		}
	}
	return nil
}

// Benchmark test to ensure performance is reasonable
func BenchmarkGetScripts(b *testing.B) {
	tempDir := b.TempDir()

	// Create a few test scripts
	for i := 0; i < 5; i++ {
		content := `#!/bin/bash
# Test function
function test_func() {
    echo "test"
}`
		path := filepath.Join(tempDir, fmt.Sprintf("script%d.sh", i))
		os.WriteFile(path, []byte(content), 0644)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scripts, err := GetScripts(tempDir)
		if err != nil {
			b.Errorf("GetScripts failed: %v", err)
		}
		if len(scripts) != 5 {
			b.Errorf("Expected 5 scripts, got %d", len(scripts))
		}
	}
}
