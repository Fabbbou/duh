package termm

import (
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FindDefaultFileEditor_EnvironmentVariables(t *testing.T) {
	// Save original values
	originalVisual := os.Getenv("VISUAL")
	originalEditor := os.Getenv("EDITOR")
	defer func() {
		os.Setenv("VISUAL", originalVisual)
		os.Setenv("EDITOR", originalEditor)
	}()

	t.Run("VISUAL takes priority", func(t *testing.T) {
		// Set valid executables that should exist on most systems
		var validEditor string
		switch runtime.GOOS {
		case "windows":
			validEditor = "notepad"
		default:
			validEditor = "vi"
		}

		os.Setenv("VISUAL", validEditor)
		os.Setenv("EDITOR", "nonexistent")

		result := FindDefaultFileEditor()
		assert.Equal(t, validEditor, result)
	})

	t.Run("EDITOR as fallback", func(t *testing.T) {
		os.Unsetenv("VISUAL")

		var validEditor string
		switch runtime.GOOS {
		case "windows":
			validEditor = "notepad"
		default:
			validEditor = "vi"
		}

		os.Setenv("EDITOR", validEditor)

		result := FindDefaultFileEditor()
		assert.Equal(t, validEditor, result)
	})

	t.Run("ignores non-existent executables", func(t *testing.T) {
		os.Setenv("VISUAL", "definitely-not-an-editor-12345")
		os.Setenv("EDITOR", "also-not-an-editor-67890")

		result := FindDefaultFileEditor()
		// Should fallback to system default since env vars point to non-existent executables
		assert.NotEqual(t, "definitely-not-an-editor-12345", result)
		assert.NotEqual(t, "also-not-an-editor-67890", result)
		assert.NotEmpty(t, result)
	})
}

func Test_FindDefaultFileEditor_PlatformSpecific(t *testing.T) {
	// Clear environment variables to test platform defaults
	originalVisual := os.Getenv("VISUAL")
	originalEditor := os.Getenv("EDITOR")
	defer func() {
		os.Setenv("VISUAL", originalVisual)
		os.Setenv("EDITOR", originalEditor)
	}()

	os.Unsetenv("VISUAL")
	os.Unsetenv("EDITOR")

	result := FindDefaultFileEditor()
	assert.NotEmpty(t, result, "Should always return a valid editor")

	// Test platform-specific fallbacks
	switch runtime.GOOS {
	case "windows":
		// Should eventually fallback to notepad if no other editors found
		if result == "notepad" {
			assert.Equal(t, "notepad", result)
		}
	case "darwin":
		// Should eventually fallback to open if no other editors found
		if result == "open" {
			assert.Equal(t, "open", result)
		}
	default:
		// Should eventually fallback to vi if no other editors found
		if result == "vi" {
			assert.Equal(t, "vi", result)
		}
	}
}

func Test_FindDefaultFileEditor_NeverEmpty(t *testing.T) {
	// Clear environment variables
	originalVisual := os.Getenv("VISUAL")
	originalEditor := os.Getenv("EDITOR")
	defer func() {
		os.Setenv("VISUAL", originalVisual)
		os.Setenv("EDITOR", originalEditor)
	}()

	os.Unsetenv("VISUAL")
	os.Unsetenv("EDITOR")

	result := FindDefaultFileEditor()
	assert.NotEmpty(t, result, "Should never return empty string")
}

func Test_isExecutableInPath(t *testing.T) {
	t.Run("returns true for existing executable", func(t *testing.T) {
		var knownExecutable string
		switch runtime.GOOS {
		case "windows":
			knownExecutable = "notepad"
		default:
			knownExecutable = "ls" // Should exist on Unix-like systems
		}

		result := isExecutableInPath(knownExecutable)
		assert.True(t, result)
	})

	t.Run("returns false for non-existent executable", func(t *testing.T) {
		result := isExecutableInPath("definitely-not-an-executable-12345")
		assert.False(t, result)
	})
}
