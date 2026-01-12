package termm

import (
	"os"
	"os/exec"
	"runtime"
)

func FindDefaultFileEditor() string {
	// First, check Unix-style environment variables
	if editor := os.Getenv("VISUAL"); editor != "" && isExecutableInPath(editor) {
		return editor
	}
	if editor := os.Getenv("EDITOR"); editor != "" && isExecutableInPath(editor) {
		return editor
	}

	// Try common editors based on platform
	var candidates []string
	switch runtime.GOOS {
	case "windows":
		candidates = []string{"code", "notepad++", "vim", "nvim", "nano", "notepad"}
	case "darwin": // macOS
		candidates = []string{"code", "vim", "nvim", "nano", "emacs", "open"}
	default: // Linux and other Unix-like
		candidates = []string{"code", "vim", "nvim", "nano", "emacs", "gedit", "kate"}
	}

	// Find first available editor
	for _, editor := range candidates {
		if isExecutableInPath(editor) {
			return editor
		}
	}

	// Platform-specific fallbacks
	switch runtime.GOOS {
	case "windows":
		return "notepad" // Always available on Windows
	case "darwin":
		return "open" // Always available on macOS
	default:
		return "vi" // Usually available on Unix-like systems
	}
}

// isExecutableInPath checks if an executable exists in the system PATH
func isExecutableInPath(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
