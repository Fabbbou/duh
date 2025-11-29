package toml_storage

import (
	"duh/internal/domain/entity"
	"errors"
	"os"
	"testing"

	"duh/internal/domain/utils"

	"github.com/stretchr/testify/assert"
)

func TestListTomlStorageRepository(t *testing.T) {
	repo, err := NewTomlStoreRepository("test_file.toml")
	assert.NoError(t, err)

	entries, err := repo.List(entity.Aliases)
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	expectedAliases := entity.StoreEntries{
		"ll":  "ls -al",
		"h-w": "echo 'Hello, World!'",
	}
	assert.Equal(t, expectedAliases, entries)

	entries, err = repo.List(entity.Exports)
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	expectedExports := entity.StoreEntries{
		"FOO": "bar",
		"BAZ": "qux",
	}
	assert.Equal(t, expectedExports, entries)

	_, err = repo.List("nonexistent-group")
	assert.Error(t, err)
	assert.Equal(t, errors.New("could not find group named nonexistent-group"), err)
}

func TestUpsertTomlStorageRepository(t *testing.T) {
	//prepare the test file
	utils.CopyFile("test_file.toml", "test_file_add_test.toml")
	defer os.Remove("test_file_add_test.toml")

	repo, err := NewTomlStoreRepository("test_file_add_test.toml")
	assert.NoError(t, err)

	err = repo.Upsert(entity.Aliases, "gs", "git status")
	assert.NoError(t, err)

	entries, err := repo.List(entity.Aliases)
	assert.NoError(t, err)
	assert.Equal(t, "git status", entries["gs"])

	err = repo.Upsert("nonexistent-group", "key", "value")
	assert.Error(t, err)
	assert.Equal(t, "group nonexistent-group does not exists", err.Error())
}

func TestDeleteTomlStorageRepository(t *testing.T) {
	utils.CopyFile("test_file.toml", "test_file_delete_test.toml")
	defer os.Remove("test_file_delete_test.toml")

	repo, err := NewTomlStoreRepository("test_file_delete_test.toml")
	assert.NoError(t, err)

	err = repo.Delete(entity.Exports, "FOO")
	assert.NoError(t, err)

	entries, err := repo.List(entity.Exports)
	assert.NoError(t, err)
	_, exists := entries["FOO"]
	assert.False(t, exists)
}
