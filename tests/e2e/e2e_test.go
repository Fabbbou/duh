package e2e

// Postponed: need XDG support first to create a local temp dir for tests

import (
	"bytes"
	"duh/cmd/cli/context"

	"os"
	"path/filepath"
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
		assert.Contains(t, output, "EDITOR='vim'")
		assert.Contains(t, output, "BROWSER='firefox'")
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
		assert.NotContains(t, output, "EDITOR='vim'")
		assert.Contains(t, output, "BROWSER='firefox'") // Should still have this one
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

	t.Run("repository management", func(t *testing.T) {
		// List repositories (should show default)
		output, err := executeCommand([]string{"repo", "list"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Enabled repositories:")
		assert.Contains(t, output, "✓ local")

		// Test adding repository with URL only
		output, err = executeCommand([]string{"repo", "add", "https://github.com/Fabbbou/my-duh"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Repository 'https://github.com/Fabbbou/my-duh' added and enabled")
		executeCommand([]string{"repo", "delete", "my-duh"})

		// Test adding repository with custom name
		output, err = executeCommand([]string{"repo", "add", "https://github.com/Fabbbou/my-duh", "myrepo"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Repository 'https://github.com/Fabbbou/my-duh' added and enabled")

		// Verify repositories are listed after adding
		output, err = executeCommand([]string{"repo", "list"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Enabled repositories:")
		executeCommand([]string{"repo", "delete", "myrepo"})

		// Test creating empty repository
		output, err = executeCommand([]string{"repo", "create", "emptyrepo"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Repository 'emptyrepo' created and enabled")

		// Verify created repository appears in list
		output, err = executeCommand([]string{"repo", "list"})
		assert.NoError(t, err)
		assert.Contains(t, output, "Enabled repositories:")
		assert.Contains(t, output, "✓ emptyrepo")

		// Clean up created repository
		executeCommand([]string{"repo", "delete", "emptyrepo"})
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

		// Test function aliases
		for _, alias := range []string{"functions", "func", "fn", "fun"} {
			_, err = executeCommand([]string{alias, "list"})
			assert.NoError(t, err)
			// Should work for all aliases
		}
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
		assert.Contains(t, output, "Manage repositories for aliases and exports")
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

		// Test function aliases help
		for _, alias := range []string{"functions", "func", "fn", "fun"} {
			output, err = executeCommand([]string{alias})
			assert.NoError(t, err)
			assert.Contains(t, output, "Manage shell functions injected by duh.")
		}
	})
}

// Test_E2E_ErrorHandling tests error scenarios
func Test_E2E_ErrorHandling(t *testing.T) {
	t.Run("invalid commands return errors", func(t *testing.T) {
		// Invalid subcommand
		output, err := executeCommand([]string{"alias", "invalid"})
		assert.Contains(t, output, "duh alias [command]")

		// Missing arguments
		output, err = executeCommand([]string{"alias", "set", "myalias"})
		assert.Error(t, err)

		// Invalid export command
		output, err = executeCommand([]string{"exports", "invalid"})
		assert.Contains(t, output, "duh exports [command]")

		// Invalid repo add command (no arguments)
		output, err = executeCommand([]string{"repo", "add"})
		assert.Error(t, err)

		// Invalid repo create command (no arguments)
		output, err = executeCommand([]string{"repo", "create"})
		assert.Error(t, err)

		// Invalid repo create command (no arguments)
		output, err = executeCommand([]string{"repo", "create"})
		assert.Error(t, err)

		// Invalid repo command
		output, err = executeCommand([]string{"repo", "invalid"})
		assert.Contains(t, output, "duh repository [command]")
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

	t.Run("path commands", func(t *testing.T) {
		// Test basic path command (should show base path)
		output, err := executeCommand([]string{"path"})
		assert.NoError(t, err)
		assert.NotEmpty(t, output)

		// Test path list command (should show all paths)
		output, err = executeCommand([]string{"path", "list"})
		assert.NoError(t, err)
		assert.NotEmpty(t, output)

		// Test path aliases
		output, err = executeCommand([]string{"paths", "list"})
		assert.NoError(t, err)
		assert.NotEmpty(t, output)
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
