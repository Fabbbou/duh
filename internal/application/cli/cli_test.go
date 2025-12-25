package cli

import (
	"bytes"
	"duh/internal/domain/entity"
	"duh/internal/domain/repository"
	"duh/internal/domain/service"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func setupMockCliService() service.CliService {
	mockRepo := &repository.MockDbRepository{
		DefaultRepo: entity.Repository{
			Name:    "default",
			Aliases: map[string]string{"ll": "ls -la"},
			Exports: map[string]string{"PATH": "/usr/bin"},
		},
		Repos: []entity.Repository{{
			Name:    "default",
			Aliases: map[string]string{"ll": "ls -la"},
			Exports: map[string]string{"PATH": "/usr/bin"},
		}},
		Enabled: []string{"default"},
	}
	return service.NewCliService(mockRepo)
}

func TestAliasCli_Help(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildAliasCli(cliService)

	// Capture output
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Execute with no args (should show help)
	cmd.SetArgs([]string{})
	err := cmd.Execute()

	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Keep alias in your shell for good, duh.")
	assert.Contains(t, output, "set")
	assert.Contains(t, output, "unset")
	assert.Contains(t, output, "list")
}

func TestAliasCli_Set(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildAliasCli(cliService)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Test setting an alias
	cmd.SetArgs([]string{"set", "gs", "git status"})
	err := cmd.Execute()

	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Alias 'gs' set for command 'git status'")
}

func TestAliasCli_List(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildAliasCli(cliService)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Test listing aliases
	cmd.SetArgs([]string{"list"})
	err := cmd.Execute()

	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "ll='ls -la'")
}

func TestAliasCli_Unset(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildAliasCli(cliService)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Test removing an alias
	cmd.SetArgs([]string{"unset", "ll"})
	err := cmd.Execute()

	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Alias 'll' removed")
}

func TestAliasCli_InvalidSubcommand(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildAliasCli(cliService)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Test invalid subcommand
	cmd.SetArgs([]string{"invalid"})
	err := cmd.Execute()

	assert.Error(t, err)
	output := buf.String()
	assert.Contains(t, output, "unknown command")
}

func TestExportsCli_Help(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildExportsCli(cliService)

	// Capture output
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Execute with no args (should show help)
	cmd.SetArgs([]string{})
	err := cmd.Execute()

	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Keep exports in your shell for good, duh.")
	assert.Contains(t, output, "set")
	assert.Contains(t, output, "unset")
	assert.Contains(t, output, "list")
}

func TestExportsCli_Set(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildExportsCli(cliService)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Test setting an export
	cmd.SetArgs([]string{"set", "EDITOR", "vim"})
	err := cmd.Execute()

	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Export 'EDITOR' set for value 'vim'")
}

func TestExportsCli_List(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildExportsCli(cliService)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Test listing exports
	cmd.SetArgs([]string{"list"})
	err := cmd.Execute()

	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "PATH='/usr/bin'")
}

func TestExportsCli_Unset(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildExportsCli(cliService)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Test removing an export
	cmd.SetArgs([]string{"unset", "PATH"})
	err := cmd.Execute()

	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Export 'PATH' removed")
}

func TestExportsCli_InvalidSubcommand(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildExportsCli(cliService)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Test invalid subcommand
	cmd.SetArgs([]string{"invalid"})
	err := cmd.Execute()

	assert.Error(t, err)
	output := buf.String()
	assert.Contains(t, output, "unknown command")
}

// Helper function to execute command and capture output
func executeCommandWithOutput(cmd *cobra.Command, args []string) (string, error) {
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs(args)

	err := cmd.Execute()
	return buf.String(), err
}

// Test the root CLI as well
func TestRootCli_Help(t *testing.T) {
	cliService := setupMockCliService()
	rootCmd := BuildRootCli(cliService)

	output, err := executeCommandWithOutput(rootCmd, []string{"--help"})

	assert.NoError(t, err)
	assert.Contains(t, output, "Duh, a simple and effective dotfiles manager")
	assert.Contains(t, output, "alias")
	assert.Contains(t, output, "exports")
	assert.Contains(t, output, "exports")
}

func TestRootCli_AliasSubcommand(t *testing.T) {
	cliService := setupMockCliService()
	rootCmd := BuildRootCli(cliService)

	output, err := executeCommandWithOutput(rootCmd, []string{"alias", "list"})

	assert.NoError(t, err)
	assert.Contains(t, output, "ll='ls -la'")
}

func TestRootCli_ExportsSubcommand(t *testing.T) {
	cliService := setupMockCliService()
	rootCmd := BuildRootCli(cliService)

	output, err := executeCommandWithOutput(rootCmd, []string{"exports", "list"})

	assert.NoError(t, err)
	assert.Contains(t, output, "PATH='/usr/bin'")
}
