package service

import (
	"duh/internal/domain/entity"
	"duh/internal/domain/repository"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectedInjection = []string{
	`alias duh_reload='eval "$(duh inject --quiet)"'`,
	`alias ll="ls -la"`,
	`alias gs="git status"`,
	`alias ca="echo \"Complex Alias\""`,
	`export PATH="/usr/local/bin:\$PATH"`,
	`export GOENV="development"`,
	`alias 2ll="2ls -la"`,
	`alias 2gs="2git status"`,
	`alias 2ca="2echo \"Complex Alias\""`,
	`export 2PATH="2/usr/local/bin:\$PATH"`,
	`export 2GOENV="2development"`,
	"",
}

var expectedInjectionStr = `alias duh_reload='eval "$(duh inject --quiet)"'
alias ll="ls -la"
alias gs="git status"
alias ca="echo \"Complex Alias\""
export PATH="/usr/local/bin:$PATH"
export GOENV="development"
alias 2ll="2ls -la"
alias 2gs="2git status"
alias 2ca="2echo \"Complex Alias\""
export 2PATH="2/usr/local/bin:$PATH"
export 2GOENV="2development"
`

func setup() CliService {
	repoDefault := entity.Repository{
		Name: "default",
		Aliases: map[string]string{
			"ll": "ls -la",
			"gs": "git status",
			"ca": `echo "Complex Alias"`,
		},
		Exports: map[string]string{
			"PATH":  "/usr/local/bin:$PATH",
			"GOENV": "development",
		},
	}

	repo2 := entity.Repository{
		Name: "second",
		Aliases: map[string]string{
			"2ll": "2ls -la",
			"2gs": "2git status",
			"2ca": `2echo "Complex Alias"`,
		},
		Exports: map[string]string{
			"2PATH":  "2/usr/local/bin:$PATH",
			"2GOENV": "2development",
		},
	}

	repo3 := entity.Repository{
		Name: "third",
		Aliases: map[string]string{
			"2ll": "2ls -la",
			"2gs": "2git status",
			"2ca": `2echo "Complex Alias"`,
		},
		Exports: map[string]string{
			"2PATH":  "2/usr/local/bin:$PATH",
			"2GOENV": "2development",
		},
	}

	mock := repository.MockDbRepository{
		DefaultRepo: repoDefault,
		Repos:       []entity.Repository{repoDefault, repo2, repo3},
		Enabled:     []string{"default", "second"},
	}

	dummyFunctionRepo := &repository.DummyFunctionRepository{}

	return NewCliService(&mock, dummyFunctionRepo)
}

func Test_Inject(t *testing.T) {
	cliService := setup()
	injection, err := cliService.Inject()
	assert.NoError(t, err)

	// Split into lines and sort for comparison
	actualLines := strings.Split(injection, "\n")
	expectedLines := strings.Split(expectedInjectionStr, "\n")

	sort.Strings(actualLines)
	sort.Strings(expectedLines)

	assert.Equal(t, expectedLines, actualLines)
}

func Test_UpsertAlias(t *testing.T) {
	cliService := setup()
	err := cliService.UpsertAlias("newalias", "newcommand")
	assert.NoError(t, err)
	aliases, err := cliService.ListAliases()
	assert.NoError(t, err)
	assert.Equal(t, "newcommand", aliases["newalias"])
}

func Test_RemoveAlias(t *testing.T) {
	cliService := setup()
	err := cliService.RemoveAlias("gs")
	assert.NoError(t, err)
	aliases, err := cliService.ListAliases()
	assert.NoError(t, err)
	_, exists := aliases["gs"]
	assert.False(t, exists)
}

func Test_ListAliases(t *testing.T) {
	cliService := setup()
	aliases, err := cliService.ListAliases()
	assert.NoError(t, err)
	expectedAliases := map[string]string{
		"ll":  "ls -la",
		"gs":  "git status",
		"ca":  `echo "Complex Alias"`,
		"2ll": "2ls -la",
		"2gs": "2git status",
		"2ca": `2echo "Complex Alias"`,
	}
	assert.Equal(t, expectedAliases, aliases)
}

func Test_SetExport(t *testing.T) {
	cliService := setup()
	err := cliService.UpsertExport("NEWEXPORT", "newvalue")
	assert.NoError(t, err)
	exports, err := cliService.ListExports()
	assert.NoError(t, err)
	assert.Equal(t, "newvalue", exports["NEWEXPORT"])
}

func Test_RemoveExport(t *testing.T) {
	cliService := setup()
	err := cliService.RemoveExport("PATH")
	assert.NoError(t, err)
	exports, err := cliService.ListExports()
	assert.NoError(t, err)
	_, exists := exports["PATH"]
	assert.False(t, exists)
}

func Test_ListExports(t *testing.T) {
	cliService := setup()
	exports, err := cliService.ListExports()
	assert.NoError(t, err)
	expectedExports := map[string]string{
		"PATH":   "/usr/local/bin:$PATH",
		"GOENV":  "development",
		"2PATH":  "2/usr/local/bin:$PATH",
		"2GOENV": "2development",
	}
	assert.Equal(t, expectedExports, exports)
}

// Repository Management Tests

func Test_ListRepositories(t *testing.T) {
	cliService := setup()
	repos, err := cliService.ListRepositories()
	assert.NoError(t, err)
	expected := map[string][]string{
		"enabled":  {"default", "second"},
		"disabled": {"third"},
	}
	assert.Equal(t, expected, repos)
}

func Test_EnableRepository(t *testing.T) {
	cliService := setup()

	err := cliService.EnableRepository("third")
	assert.NoError(t, err)
	repos, err := cliService.ListRepositories()
	assert.NoError(t, err)
	assert.Contains(t, repos["enabled"], "third")
	assert.NotContains(t, repos["disabled"], "third")
}

func Test_DisableRepository(t *testing.T) {
	cliService := setup()
	err := cliService.DisableRepository("second")
	assert.NoError(t, err)
	repos, err := cliService.ListRepositories()
	assert.NoError(t, err)
	assert.Contains(t, repos["disabled"], "second")
	assert.NotContains(t, repos["enabled"], "second")
}

func Test_DeleteRepository(t *testing.T) {
	cliService := setup()
	err := cliService.DeleteRepository("second")
	assert.NoError(t, err)
	repos, err := cliService.ListRepositories()
	assert.NoError(t, err)
	assert.NotContains(t, repos["enabled"], "second")
	assert.NotContains(t, repos["disabled"], "second")
}

func Test_SetDefaultRepository(t *testing.T) {
	cliService := setup()
	err := cliService.SetDefaultRepository("third")
	assert.NoError(t, err)
	repos, err := cliService.ListRepositories()
	assert.NoError(t, err)
	assert.Contains(t, repos["enabled"], "third")
}

func Test_RenameRepository(t *testing.T) {
	cliService := setup()
	err := cliService.RenameRepository("second", "renamed")
	assert.NoError(t, err)
	repos, err := cliService.ListRepositories()
	assert.NoError(t, err)
	assert.Contains(t, repos["enabled"], "renamed")
	assert.NotContains(t, repos["enabled"], "second")
}

func Test_GetCurrentDefaultRepository(t *testing.T) {
	cliService := setup()
	defaultRepo, err := cliService.GetCurrentDefaultRepository()
	assert.NoError(t, err)
	assert.Equal(t, "default", defaultRepo)
}

func Test_AddRepository(t *testing.T) {
	cliService := setup()

	// Test adding repository with URL only
	err := cliService.AddRepository("https://github.com/Fabbbou/my-duh", nil)
	assert.NoError(t, err)

	// Test adding repository with custom name
	name := "customrepo"
	err = cliService.AddRepository("https://github.com/Fabbbou/my-duh", &name)
	assert.NoError(t, err)
}

func Test_CreateRepository(t *testing.T) {
	cliService := setup()

	// Test creating repository
	err := cliService.CreateRepository("newrepo")
	assert.NoError(t, err)
	repos, err := cliService.ListRepositories()
	assert.NoError(t, err)
	assert.Contains(t, repos["enabled"], "newrepo")
}

func Test_PushRepository(t *testing.T) {
	cliService := setup()

	// Test pushing repository (will succeed with mock)
	err := cliService.PushRepository("default")
	assert.NoError(t, err)
}
