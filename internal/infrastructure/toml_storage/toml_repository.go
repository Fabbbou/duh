package toml_storage

import (
	"cmd/cli/main.go/internal/domain/entity"
	"errors"
)

type TomlStorageRepository struct {
	filepath string
	storage  entity.StoreGroup
}

func NewTomlStorageRepository(filepath string) (*TomlStorageRepository, error) {
	storage, err := loadFile(filepath)
	if err != nil {
		return nil, err
	}

	return &TomlStorageRepository{
		filepath: filepath,
		storage:  map_storage_to_entity(storage),
	}, nil
}

func (r *TomlStorageRepository) Upsert(groupName entity.GroupName, key entity.Key, value entity.Value) error {
	// Implementation for upserting a config entry into the TOML file
	return nil
}

func (r *TomlStorageRepository) List(groupName entity.GroupName) ([]entity.StorageEntry, error) {
	entries, exists := r.storage[groupName]
	if !exists {
		return nil, errors.New("could not find group named " + string(groupName))
	}
	return entries, nil
}

func (r *TomlStorageRepository) Delete(groupName entity.GroupName, key entity.Key) error {
	// Implementation for deleting a config entry from the TOML file
	return nil
}
