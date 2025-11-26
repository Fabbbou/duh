package toml_storage

import (
	"cmd/cli/main.go/internal/domain/entity"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListTomlStorageRepository(t *testing.T) {
	repo, err := NewTomlStorageRepository("test_file.toml")
	assert.NoError(t, err)

	entries, err := repo.List(entity.Aliases)
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	assert.Equal(t, entity.StorageEntry{Key: "ll", Value: "ls -al"}, entries[0])
	assert.Equal(t, entity.StorageEntry{Key: "h-w", Value: "echo 'Hello, World!'"}, entries[1])

	entries, err = repo.List(entity.Exports)
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	assert.Equal(t, entity.StorageEntry{Key: "FOO", Value: "bar"}, entries[0])
	assert.Equal(t, entity.StorageEntry{Key: "BAZ", Value: "qux"}, entries[1])

	_, err = repo.List("nonexistent-group")
	assert.Error(t, err)
	assert.Equal(t, errors.New("could not find group named nonexistent-group"), err)
}
