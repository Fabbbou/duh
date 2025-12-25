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
	`alias duh_inject="eval $(duh inject)"`,
	`alias ll="ls -la"`,
	`alias gs="git status"`,
	`alias ca="echo \"Complex Alias\""`,
	`export PATH="/usr/local/bin:$PATH"`,
	`export GOENV="development"`,
	`alias 2ll="2ls -la"`,
	`alias 2gs="2git status"`,
	`alias 2ca="2echo \"Complex Alias\""`,
	`export 2PATH="2/usr/local/bin:$PATH"`,
	`export 2GOENV="2development"`,
}

var expectedInjectionStr = `alias duh_inject="eval $(duh inject)"
alias ll="ls -la"
alias gs="git status"
alias ca="echo \"Complex Alias\""
export PATH="/usr/local/bin:$PATH"
export GOENV="development"
alias 2ll="2ls -la"
alias 2gs="2git status"
alias 2ca="2echo \"Complex Alias\""
export 2PATH="2/usr/local/bin:$PATH"
export 2GOENV="2development"`

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

	mock := repository.MockDbRepository{
		DefaultRepo: repoDefault,
		Repos:       []entity.Repository{repoDefault, repo2},
		Enabled:     []string{"default", "second"},
	}

	return NewCliService(&mock)
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
