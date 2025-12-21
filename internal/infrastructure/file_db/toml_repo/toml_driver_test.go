package toml_repo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDriverLoad(t *testing.T) {
	storage, err := LoadToml[RepositoryToml]("test_file.toml")
	if err != nil {
		t.FailNow()
	}

	assert.NoError(t, err)

	assert.Len(t, storage.Aliases, 2)
	assert.Equal(t, "ls -al", storage.Aliases["ll"])
	assert.Equal(t, "echo 'Hello, World!'", storage.Aliases["h-w"])

	assert.Len(t, storage.Exports, 2)
	assert.Equal(t, "bar", storage.Exports["FOO"])
	assert.Equal(t, "qux", storage.Exports["BAZ"])
}

func TestDriverSave(t *testing.T) {

	createStorage := &RepositoryToml{
		Aliases: map[string]string{
			"gs": "git status",
		},
		Exports: map[string]string{
			"PATH": "/usr/local/bin",
		},
	}

	err := SaveToml("test_save.toml", *createStorage)
	assert.NoError(t, err)

	storage, err := LoadToml[RepositoryToml]("test_save.toml")
	if err != nil {
		t.FailNow()
	}
	assert.Len(t, storage.Aliases, 1)
	assert.Equal(t, "git status", storage.Aliases["gs"])

	assert.Len(t, storage.Exports, 1)
	assert.Equal(t, "/usr/local/bin", storage.Exports["PATH"])

	os.Remove("test_save.toml")
}
