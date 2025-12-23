package service

import (
	"duh/internal/domain/entity"
	"duh/internal/domain/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectedInjection = []string{
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

var expectedInjectionStr = `alias ll="ls -la"
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
	mock := repository.MockDbRepository{
		DefaultRepo: repoDefault,
		Repos:       []entity.Repository{repoDefault},
		Enabled:     []string{"default"},
	}

	return NewCliService(&mock)
}

func Test_Inject(t *testing.T) {
	cliService := setup()
	injection, err := cliService.Inject()
	assert.NoError(t, err)
	assert.Equal(t, expectedInjectionStr, injection)
}
