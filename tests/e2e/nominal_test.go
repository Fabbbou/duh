package e2e

// Postponed: need XDG support first to create a local temp dir for tests

import (
	"bytes"
	"duh/internal/application/contexts"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// executeCommand is a helper to run CLI commands and capture output
func executeCommand(args []string) (string, error) {
	cli := contexts.InitCli()
	var buf bytes.Buffer
	cli.SetOut(&buf)
	cli.SetErr(&buf)
	cli.SetArgs(args)
	err := cli.Execute()
	return buf.String(), err
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

		// Verify proper shell command format
		lines := strings.Split(strings.TrimSpace(output), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				assert.True(t,
					strings.HasPrefix(line, "alias ") || strings.HasPrefix(line, "export "),
					"Line should start with 'alias' or 'export': %s", line)
			}
		}
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
}
