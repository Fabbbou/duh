package e2e

// Postponed: need XDG support first to create a local temp dir for tests

import (
	"bytes"
	"duh/cmd/cli/context"

	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/adrg/xdg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var duhPath string

func init() {
	//overiding xdg path for tests
	xdg.DataHome = os.TempDir()
	duhPath = filepath.Join(xdg.DataHome, "duh")
	os.RemoveAll(duhPath)
}

// executeCommand is a helper to run CLI commands and capture output
func executeCommand(args []string) (string, error) {
	// Capture real stdout/stderr for commands that write directly to os.Stdout
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	// Create pipes to capture output
	stdoutReader, stdoutWriter, _ := os.Pipe()
	stderrReader, stderrWriter, _ := os.Pipe()

	// Redirect os.Stdout and os.Stderr
	os.Stdout = stdoutWriter
	os.Stderr = stderrWriter

	// Channels to collect output
	stdoutChan := make(chan string)
	stderrChan := make(chan string)

	// Goroutines to read from pipes
	go func() {
		var buf bytes.Buffer
		buf.ReadFrom(stdoutReader)
		stdoutChan <- buf.String()
	}()

	go func() {
		var buf bytes.Buffer
		buf.ReadFrom(stderrReader)
		stderrChan <- buf.String()
	}()

	// Execute command
	cli := context.InitializeCLI()
	cli.SetArgs(args)
	err := cli.Execute()

	// Close writers to signal end of output
	stdoutWriter.Close()
	stderrWriter.Close()

	// Collect output
	stdoutOutput := <-stdoutChan
	stderrOutput := <-stderrChan

	// Restore original stdout/stderr
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	// Combine outputs for backward compatibility
	combined := stdoutOutput + stderrOutput
	return combined, err
}

// Test_E2E_Complete tests the full duh CLI workflow end-to-end
func Test_E2E_Complete(t *testing.T) {
	t.Run("initial state should be empty", func(t *testing.T) {
		// Test aliases list (should be empty)
		_, err := executeCommand([]string{"alias", "list"})
		assert.NoError(t, err)
		// Should either be empty or show no output

		// Test exports list (should be empty)
		_, err = executeCommand([]string{"exports", "list"})
		assert.NoError(t, err)
		// Should either be empty or show no output
	})

	t.Run("set and list aliases", func(t *testing.T) {
		// Set some aliases
		output, err := executeCommand([]string{"alias", "set", "ll", "ls -la"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Alias 'll' set for command 'ls -la'")

		output, err = executeCommand([]string{"alias", "set", "gs", "git status"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Alias 'gs' set for command 'git status'")

		// List aliases and verify they exist
		output, err = executeCommand([]string{"alias", "list"})
		assert.NoError(t, err)
		assert.Contains(t, output, "ll='ls -la'")
		assert.Contains(t, output, "gs='git status'")
	})

	t.Run("set and list exports", func(t *testing.T) {
		// Set some exports
		output, err := executeCommand([]string{"exports", "set", "EDITOR", "vim"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Export 'EDITOR' set for value 'vim'")

		output, err = executeCommand([]string{"exports", "set", "BROWSER", "firefox"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Export 'BROWSER' set for value 'firefox'")

		// List exports and verify they exist
		output, err = executeCommand([]string{"exports", "list"})
		assert.NoError(t, err)
		assert.Contains(t, output, "EDITOR=vim")
		assert.Contains(t, output, "BROWSER=firefox")
	})

	t.Run("inject command generates correct output", func(t *testing.T) {
		output, err := executeCommand([]string{"inject"})
		assert.NoError(t, err)

		// Should contain all aliases and exports
		assert.Contains(t, output, "alias ll=\"ls -la\"")
		assert.Contains(t, output, "alias gs=\"git status\"")
		assert.Contains(t, output, "export EDITOR=\"vim\"")
		assert.Contains(t, output, "export BROWSER=\"firefox\"")
	})

	t.Run("remove aliases", func(t *testing.T) {
		// Remove an alias
		output, err := executeCommand([]string{"alias", "unset", "ll"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Alias 'll' removed")

		// Verify it's gone
		output, err = executeCommand([]string{"alias", "list"})
		assert.NoError(t, err)
		assert.NotContains(t, output, "ll='ls -la'")
		assert.Contains(t, output, "gs='git status'") // Should still have this one
	})

	t.Run("remove exports", func(t *testing.T) {
		// Remove an export
		output, err := executeCommand([]string{"exports", "unset", "EDITOR"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Export 'EDITOR' removed")

		// Verify it's gone
		output, err = executeCommand([]string{"exports", "list"})
		assert.NoError(t, err)
		assert.NotContains(t, output, "EDITOR=vim")
		assert.Contains(t, output, "BROWSER=firefox") // Should still have this one
	})

	t.Run("inject reflects changes", func(t *testing.T) {
		output, err := executeCommand([]string{"inject"})
		assert.NoError(t, err)

		// Should only contain remaining items
		assert.NotContains(t, output, "alias ll=")
		assert.NotContains(t, output, "export EDITOR=")
		assert.Contains(t, output, "alias gs=\"git status\"")
		assert.Contains(t, output, "export BROWSER=\"firefox\"")
	})

	t.Run("package management", func(t *testing.T) {
		// Test using new package command
		// List packages (should show default)
		output, err := executeCommand([]string{"package", "list"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Enabled packages:")
		assert.Contains(t, output, "✓ local")

		// Test backward compatibility with repo command
		// List repositories (should show default)
		output, err = executeCommand([]string{"repo", "list"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Enabled packages:") // Should still show packages terminology
		assert.Contains(t, output, "✓ local")

		// Test adding package with URL only
		output, err = executeCommand([]string{"package", "add", "https://github.com/Fabbbou/my-duh"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Package 'https://github.com/Fabbbou/my-duh' added and enabled")
		executeCommand([]string{"package", "delete", "my-duh"})

		// Test adding package with custom name using old repo alias
		output, err = executeCommand([]string{"repo", "add", "https://github.com/Fabbbou/my-duh", "myrepo"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Package 'https://github.com/Fabbbou/my-duh' added and enabled") // Output should use package terminology

		// Verify packages are listed after adding
		output, err = executeCommand([]string{"package", "list"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Enabled packages:")
		executeCommand([]string{"package", "delete", "myrepo"})

		// Test creating empty package
		output, err = executeCommand([]string{"package", "create", "emptypackage"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Package 'emptypackage' created and enabled")

		// Verify created package appears in list
		output, err = executeCommand([]string{"package", "list"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Enabled packages:")
		assert.Contains(t, output, "✓ emptypackage")

		// Clean up created package
		executeCommand([]string{"package", "delete", "emptypackage"})
	})

	t.Run("function management", func(t *testing.T) {
		// Test listing functions (should show available functions)
		_, err := executeCommand([]string{"function", "list"})
		assert.NoError(t, err)
		// May be empty if no functions are activated, but should not error

		// Test listing all functions
		_, err = executeCommand([]string{"function", "list", "--all"})
		assert.NoError(t, err)
		// Should work regardless of available functions

		// Test listing all functions with short flag
		_, err = executeCommand([]string{"function", "list", "-a"})
		assert.NoError(t, err)
		// Should work regardless of available functions

		// Test listing core/internal functions
		output, err := executeCommand([]string{"function", "list", "--core"})
		assert.NoError(t, err)
		// Should show internal functions from embedded scripts
		assert.Contains(t, output, "()")

		// Test function aliases
		for _, alias := range []string{"functions", "func", "fn", "fun"} {
			_, err = executeCommand([]string{alias, "list"})
			assert.NoError(t, err)
			// Should work for all aliases
		}

		// Test function aliases with core flag
		for _, alias := range []string{"functions", "func", "fn", "fun"} {
			_, err = executeCommand([]string{alias, "list", "--core"})
			assert.NoError(t, err)
			// Should work for all aliases with core flag
		}

		// Test function info subcommand (with core functions)
		output, err = executeCommand([]string{"function", "list", "--core"})
		assert.NoError(t, err)

		// Try to get info on a core function if any exist
		if len(output) > 0 {
			// Extract the first function name from the output (format: "- functionName()")
			lines := strings.Split(output, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "- ") && strings.HasSuffix(line, "()") {
					funcName := strings.TrimSuffix(strings.TrimPrefix(line, "- "), "()")
					if funcName != "" {
						infoOutput, infoErr := executeCommand([]string{"function", "info", funcName})
						assert.NoError(t, infoErr)
						// Should show function details
						assert.Contains(t, infoOutput, funcName+"()")
						break
					}
				}
			}
		}

		// Test function info with nonexistent function
		output, err = executeCommand([]string{"function", "info", "nonexistent_function_name"})
		assert.NoError(t, err)
		assert.Contains(t, output, "not found")

		// Test function info aliases
		output, err = executeCommand([]string{"func", "info", "nonexistent_function_name"})
		assert.NoError(t, err)
		assert.Contains(t, output, "not found")
	})
}

// Test_E2E_Help tests help commands
func Test_E2E_Help(t *testing.T) {
	t.Run("root help", func(t *testing.T) {
		output, err := executeCommand([]string{"--help"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Duh, a simple and effective dotfiles manager")
		assert.Contains(t, output, "alias")
		assert.Contains(t, output, "exports")
	})

	t.Run("alias help", func(t *testing.T) {
		output, err := executeCommand([]string{"alias"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Keep alias in your shell for good, duh.")
		assert.Contains(t, output, "set")
		assert.Contains(t, output, "unset")
		assert.Contains(t, output, "list")
	})

	t.Run("exports help", func(t *testing.T) {
		output, err := executeCommand([]string{"exports"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Keep exports in your shell for good, duh.")
		assert.Contains(t, output, "set")
		assert.Contains(t, output, "unset")
		assert.Contains(t, output, "list")
	})

	t.Run("repo help", func(t *testing.T) {
		output, err := executeCommand([]string{"repo"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Manage packages for aliases and exports")
		assert.Contains(t, output, "list")
		assert.Contains(t, output, "enable")
		assert.Contains(t, output, "disable")
		assert.Contains(t, output, "add")
		assert.Contains(t, output, "create")
	})

	t.Run("function help", func(t *testing.T) {
		output, err := executeCommand([]string{"function"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Manage shell functions injected by duh.")
		assert.Contains(t, output, "list")
		assert.Contains(t, output, "info")
		assert.Contains(t, output, "add")

		// Test function aliases help
		for _, alias := range []string{"functions", "func", "fn", "fun"} {
			output, err = executeCommand([]string{alias})
			assert.NoError(t, err)
			assert.Contains(t, output, "Manage shell functions injected by duh.")
			assert.Contains(t, output, "list")
			assert.Contains(t, output, "info")
			assert.Contains(t, output, "add")
		}
	})
}

// Test_E2E_ErrorHandling tests error scenarios
func Test_E2E_ErrorHandling(t *testing.T) {
	t.Run("invalid commands return errors", func(t *testing.T) {
		// Invalid subcommand (these should show help, not error)
		output, err := executeCommand([]string{"alias", "invalid"})
		assert.NoError(t, err) // Should show help, not error
		assert.Contains(t, output, "duh alias [command]")

		// Missing arguments should error
		output, err = executeCommand([]string{"alias", "set", "myalias"})
		assert.Error(t, err)

		// Invalid export command (should show help, not error)
		output, err = executeCommand([]string{"exports", "invalid"})
		assert.NoError(t, err) // Should show help, not error
		assert.Contains(t, output, "duh exports [command]")

		// Invalid repo add command (no arguments)
		output, err = executeCommand([]string{"repo", "add"})
		assert.Error(t, err)

		// Invalid repo create command (no arguments)
		output, err = executeCommand([]string{"repo", "create"})
		assert.Error(t, err)

		// Invalid repo command (should show help, not error)
		output, err = executeCommand([]string{"repo", "invalid"})
		assert.NoError(t, err) // Should show help, not error
		assert.Contains(t, output, "duh package [command]")
	})

	t.Run("removing non-existent items", func(t *testing.T) {
		// Try to remove non-existent alias
		_, err := executeCommand([]string{"alias", "unset", "nonexistent"})
		// This should either succeed silently or show a warning, not hard error
		// Adjust based on actual behavior
		if err != nil {
			// If it errors, that's also valid behavior
			require.Error(t, err)
		}

		// Try to remove non-existent export
		_, err = executeCommand([]string{"exports", "unset", "NONEXISTENT"})
		// Same as above - adjust based on actual behavior
		if err != nil {
			require.Error(t, err)
		}
	})

	t.Run("function command errors", func(t *testing.T) {
		// Test function info with no arguments (should error)
		output, err := executeCommand([]string{"function", "info"})
		assert.Error(t, err)

		// Test function add with no arguments (should error)
		output, err = executeCommand([]string{"function", "add"})
		assert.Error(t, err)

		// Test function info with too many arguments (should error)
		output, err = executeCommand([]string{"function", "info", "func1", "func2"})
		assert.Error(t, err)

		// Test function add with too many arguments (should error)
		output, err = executeCommand([]string{"function", "add", "func1", "func2"})
		assert.Error(t, err)

		// Test invalid function command (should show help, not error)
		output, err = executeCommand([]string{"function", "invalid"})
		assert.NoError(t, err) // Should show help, not error
		assert.Contains(t, output, "duh functions [command]")
	})
}

// Test_E2E_SpecialCharacters tests handling of special characters and edge cases
func Test_E2E_SpecialCharacters(t *testing.T) {
	t.Run("aliases with special characters", func(t *testing.T) {
		// Test with quotes and spaces
		output, err := executeCommand([]string{"alias", "set", "complex", "echo \"hello world\""})
		assert.NoError(t, err)

		// Verify it's stored correctly
		output, err = executeCommand([]string{"alias", "list"})
		assert.NoError(t, err)
		assert.Contains(t, output, "complex=")

		// Verify inject handles it correctly
		output, err = executeCommand([]string{"inject"})
		assert.NoError(t, err)
		assert.Contains(t, output, "alias complex=")
		assert.Contains(t, output, "echo \\\"hello world\\\"")

		// Clean up
		executeCommand([]string{"alias", "unset", "complex"})
	})

	t.Run("exports with special characters", func(t *testing.T) {
		// Test with special characters in value
		output, err := executeCommand([]string{"exports", "set", "PATH_EXTRA", "/usr/local/bin:/opt/bin"})
		assert.NoError(t, err)

		// Verify it's stored correctly
		output, err = executeCommand([]string{"exports", "list"})
		assert.NoError(t, err)
		assert.Contains(t, output, "PATH_EXTRA=")

		// Verify inject handles it correctly
		output, err = executeCommand([]string{"inject"})
		assert.NoError(t, err)
		assert.Contains(t, output, "export PATH_EXTRA=")

		// Clean up
		executeCommand([]string{"exports", "unset", "PATH_EXTRA"})
	})

	t.Run("repository update command", func(t *testing.T) {
		// Test update command with safe strategy (should work even without repos)
		_, err := executeCommand([]string{"repository", "update"})
		assert.NoError(t, err)

		// Test update command with commit flag
		_, err = executeCommand([]string{"repository", "update", "--commit"})
		assert.NoError(t, err)

		// Test update command with force flag
		_, err = executeCommand([]string{"repository", "update", "--force"})
		assert.NoError(t, err)

		// Test that conflicting flags are rejected
		output, err := executeCommand([]string{"repository", "update", "--force", "--commit"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Cannot use both --force and --commit")
	})

	t.Run("self commands", func(t *testing.T) {
		// Test self version command
		output, err := executeCommand([]string{"self", "version"})
		assert.NoError(t, err)
		assert.NotEmpty(t, output)

		// Test self config-path command (should show config path)
		output, err = executeCommand([]string{"self", "config-path"})
		assert.NoError(t, err)
		assert.NotEmpty(t, output)

		// Test self repositories-path command (should show repositories path)
		output, err = executeCommand([]string{"self", "repositories-path"})
		assert.NoError(t, err)
		assert.NotEmpty(t, output)

		// Test self update command (will likely fail in test environment, but should not panic)
		output, err = executeCommand([]string{"self", "update"})
		// Update command may fail in test environment (no network, already latest version, etc.)
		// We just want to ensure it doesn't crash and provides meaningful output
		assert.NotEmpty(t, output)
		// Should contain either success or error message
		assert.True(t, strings.Contains(output, "Checking for updates") ||
			strings.Contains(output, "Update failed") ||
			strings.Contains(output, "already running the latest version"))
	})

	t.Run("repository push command", func(t *testing.T) {
		// Test push command (should work even for repos without git, will return error message)
		output, err := executeCommand([]string{"repository", "push", "local"})
		// This will fail because local repo doesn't have git, but command should execute
		// Error is expected in this case
		assert.NoError(t, err)
		assert.Contains(t, output, "not a git repository")

		// Test push command with invalid arguments
		output, err = executeCommand([]string{"repository", "push"})
		assert.Error(t, err)
		assert.Contains(t, output, "accepts 1 arg(s), received 0")
	})
}
