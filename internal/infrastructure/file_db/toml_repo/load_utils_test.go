package toml_repo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadRepository(t *testing.T) {
	repo, err := LoadRepository("test_file.toml")
	assert.NoError(t, err)
	assert.Equal(t, RepositoryToml{
		Aliases: map[string]string{
			"ll":  "ls -al",
			"h-w": "echo 'Hello, World!'",
		},
		Exports: map[string]string{
			"FOO": "bar",
			"BAZ": "qux",
		},
		Metadata: MetadataMap{
			UrlOrigin:  "https://github.com/Fabbbou/duh-test-repo",
			NameOrigin: "duh-test-repo",
		},
	}, *repo)
}
