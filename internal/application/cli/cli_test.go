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
	cmd := BuildAliasSubcommand(cliService)

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
	cmd := BuildAliasSubcommand(cliService)

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
	cmd := BuildAliasSubcommand(cliService)

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
	cmd := BuildAliasSubcommand(cliService)

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
	cmd := BuildAliasSubcommand(cliService)

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
	cmd := BuildExportsSubcommand(cliService)

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
	cmd := BuildExportsSubcommand(cliService)

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
	cmd := BuildExportsSubcommand(cliService)

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
	cmd := BuildExportsSubcommand(cliService)

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
	cmd := BuildExportsSubcommand(cliService)

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
	assert.Contains(t, output, "repo")
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

// Test Repo CLI
func TestRepoCli_Help(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildRepoSubcommand(cliService)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	cmd.SetArgs([]string{})
	err := cmd.Execute()

	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Manage repositories for aliases and exports")
	assert.Contains(t, output, "list")
	assert.Contains(t, output, "enable")
	assert.Contains(t, output, "disable")
}

func TestRepoCli_List(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildRepoSubcommand(cliService)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	cmd.SetArgs([]string{"list"})
	err := cmd.Execute()

	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Enabled repositories:")
	assert.Contains(t, output, "✓ default")
}

func TestRepoCli_Enable(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildRepoSubcommand(cliService)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	cmd.SetArgs([]string{"enable", "testrepo"})
	err := cmd.Execute()

	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Repository 'testrepo' enabled")
}

func TestRepoCli_Disable(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildRepoSubcommand(cliService)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	cmd.SetArgs([]string{"disable", "default"})
	err := cmd.Execute()

	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Repository 'default' disabled")
}

func TestRepoCli_Default(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildRepoSubcommand(cliService)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	cmd.SetArgs([]string{"default"})
	err := cmd.Execute()

	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Current default repository:")
	assert.Contains(t, output, "Available commands:")
}

func TestRepoCli_DefaultSet(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildRepoSubcommand(cliService)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	cmd.SetArgs([]string{"default", "set", "newdefault"})
	err := cmd.Execute()

	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Repository 'newdefault' set as default")
}

func TestRepoCli_Add(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildRepoSubcommand(cliService)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Test adding repository with URL only
	cmd.SetArgs([]string{"add", "https://github.com/user/repo.git"})
	err := cmd.Execute()

	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Repository 'https://github.com/user/repo.git' added and enabled")
}

func TestRepoCli_AddWithName(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildRepoSubcommand(cliService)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Test adding repository with URL and custom name
	cmd.SetArgs([]string{"add", "https://github.com/user/repo.git", "myrepo"})
	err := cmd.Execute()

	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Repository 'https://github.com/user/repo.git' added and enabled")
}

func TestRepoCli_AddInvalidArgs(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildRepoSubcommand(cliService)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Test adding repository with no arguments (should fail)
	cmd.SetArgs([]string{"add"})
	err := cmd.Execute()

	assert.Error(t, err)
}

func TestRepoCli_Create(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildRepoSubcommand(cliService)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Test creating repository
	cmd.SetArgs([]string{"create", "newrepo"})
	err := cmd.Execute()

	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Repository 'newrepo' created and enabled")
}

func TestRepoCli_CreateInvalidArgs(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildRepoSubcommand(cliService)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Test creating repository with no arguments (should fail)
	cmd.SetArgs([]string{"create"})
	err := cmd.Execute()

	assert.Error(t, err)
}

func TestRootCli_RepoSubcommand(t *testing.T) {
	cliService := setupMockCliService()
	rootCmd := BuildRootCli(cliService)

	output, err := executeCommandWithOutput(rootCmd, []string{"repo", "list"})

	assert.NoError(t, err)
	assert.Contains(t, output, "Enabled repositories:")
}

func TestPathCli_Help(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildPathSubcommand(cliService)

	// Capture output
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Execute with help flag
	cmd.SetArgs([]string{"--help"})
	err := cmd.Execute()

	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Manage and view repository paths")
	assert.Contains(t, output, "list")
}

func TestPathCli_ShowBasePath(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildPathSubcommand(cliService)

	_, err := executeCommandWithOutput(cmd, []string{})

	assert.NoError(t, err)

	// Should return the base path from the mock
}

func TestPathCli_ListPaths(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildPathSubcommand(cliService)

	_, err := executeCommandWithOutput(cmd, []string{"list"})

	assert.NoError(t, err)

	// Should include both base path and repository paths
}

func TestRootCli_PathSubcommand(t *testing.T) {
	cliService := setupMockCliService()
	rootCmd := BuildRootCli(cliService)

	_, err := executeCommandWithOutput(rootCmd, []string{"path", "list"})

	assert.NoError(t, err)
}

func TestRepoCli_Push(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildRepoSubcommand(cliService)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Test pushing a repository (will succeed with mock)
	cmd.SetArgs([]string{"push", "default"})
	err := cmd.Execute()

	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Repository 'default' pushed successfully")
}

func TestRepoCli_PushInvalidArgs(t *testing.T) {
	cliService := setupMockCliService()
	cmd := BuildRepoSubcommand(cliService)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Test pushing repository with no arguments (should fail)
	cmd.SetArgs([]string{"push"})
	err := cmd.Execute()

	assert.Error(t, err)
}
