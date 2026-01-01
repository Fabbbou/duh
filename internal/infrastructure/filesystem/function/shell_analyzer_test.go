package function

import (
	"testing"
)

func TestShellAnalyzer_AnalyzeScript(t *testing.T) {
	testScript := `#!/bin/bash

# This is some global code that shouldn't be here
echo "Global execution - bad practice!"

# a comment that should be ignored

# This function handles user authentication
# It takes username and password as parameters
# Returns 0 on success, 1 on failure
function authenticate_user() {
    local username="$1"
    local password="$2"
    
    if [ "$username" = "admin" ] && [ "$password" = "secret" ]; then
        return 0
    else
        return 1
    fi
}

# Another function without proper documentation
function cleanup() {
    rm -rf /tmp/myapp/*
}

# Process data files
process_files() {
    for file in "$@"; do
        echo "Processing $file"
    done
}

# More global code - also bad
export GLOBAL_VAR="should be in function"

# Main function that orchestrates the application
# Usage: main [options]
main() {
    authenticate_user "admin" "secret"
    cleanup
    process_files *.txt
}

# Call main - this is global execution too
main "$@"`

	analyzer := NewShellAnalyzer()
	err := analyzer.AnalyzeScript(testScript)
	if err != nil {
		t.Fatalf("Failed to analyze script: %v", err)
	}

	// Check functions found
	if len(analyzer.Functions) != 4 {
		t.Errorf("Expected 4 functions, got %d", len(analyzer.Functions))
	}

	// Check function with documentation
	foundAuthFunction := false
	foundMainFunction := false
	for _, fn := range analyzer.Functions {
		if fn.Name == "authenticate_user" {
			foundAuthFunction = true
			if !fn.HasDocs {
				t.Errorf("authenticate_user should have documentation")
			}
			if len(fn.Documentation) != 3 {
				t.Errorf("authenticate_user should have 3 documentation lines, got %d", len(fn.Documentation))
			}
		}
		if fn.Name == "main" {
			foundMainFunction = true
			if !fn.HasDocs {
				t.Errorf("main function should have documentation")
			}
		}
		if fn.Name == "cleanup" {
			if !fn.HasDocs {
				t.Errorf("cleanup function should have documentation (comment above it)")
			}
			// The comment "Another function without proper documentation"
			// is still documentation, even if it's not helpful
		}
	}

	if !foundAuthFunction {
		t.Error("Should have found authenticate_user function")
	}
	if !foundMainFunction {
		t.Error("Should have found main function")
	}

	// Check code outside functions
	if len(analyzer.CodeOutside) == 0 {
		t.Error("Should have found code outside functions")
	}

	// Print the report for manual verification
	t.Log("Analysis Report:")
	analyzer.PrintReport()
}

func TestShellAnalyzer_CleanScript(t *testing.T) {
	cleanScript := `#!/bin/bash

# Utility function to check if file exists
# Parameters: filepath
# Returns: 0 if exists, 1 if not
file_exists() {
    [ -f "$1" ]
}

# Main application entry point
# Handles all application logic
main() {
    if file_exists "config.txt"; then
        echo "Config found"
    else
        echo "Config missing"
    fi
}`

	analyzer := NewShellAnalyzer()
	err := analyzer.AnalyzeScript(cleanScript)
	if err != nil {
		t.Fatalf("Failed to analyze clean script: %v", err)
	}

	// Should find 2 functions
	if len(analyzer.Functions) != 2 {
		t.Errorf("Expected 2 functions, got %d", len(analyzer.Functions))
	}

	// Should find no code outside functions
	if len(analyzer.CodeOutside) != 0 {
		t.Errorf("Expected no code outside functions, got %d", len(analyzer.CodeOutside))
	}

	// Both functions should have documentation
	for _, fn := range analyzer.Functions {
		if !fn.HasDocs {
			t.Errorf("Function %s should have documentation", fn.Name)
		}
	}

	t.Log("Clean script analysis:")
	analyzer.PrintReport()
}
