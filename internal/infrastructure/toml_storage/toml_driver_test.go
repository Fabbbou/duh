package toml_storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDriverLoad(t *testing.T) {
	driver := TomlDriver{
		filePath: "test_file.toml",
	}
	storage, err := driver.Load()

	assert.NoError(t, err)

	assert.Len(t, storage.Aliases, 2)
	assert.Equal(t, "ls -al", storage.Aliases["ll"])
	assert.Equal(t, "echo 'Hello, World!'", storage.Aliases["h-w"])

	assert.Len(t, storage.Exports, 2)
	assert.Equal(t, "bar", storage.Exports["FOO"])
	assert.Equal(t, "qux", storage.Exports["BAZ"])
}

func TestDriverSave(t *testing.T) {
	driver := TomlDriver{
		filePath: "test_save.toml",
	}
	newStorage := &Storage{
		Aliases: map[string]string{
			"gs": "git status",
		},
		Exports: map[string]string{
			"PATH": "/usr/local/bin",
		},
	}

	err := driver.Save(newStorage)
	assert.NoError(t, err)

	loadedStorage, err := driver.Load()
	assert.NoError(t, err)

	assert.Len(t, loadedStorage.Aliases, 1)
	assert.Equal(t, "git status", loadedStorage.Aliases["gs"])

	assert.Len(t, loadedStorage.Exports, 1)
	assert.Equal(t, "/usr/local/bin", loadedStorage.Exports["PATH"])

	os.Remove("test_save.toml")
}
