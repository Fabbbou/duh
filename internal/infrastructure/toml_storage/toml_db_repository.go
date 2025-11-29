package toml_storage

import (
	"duh/internal/domain/entity"
	"errors"
	"fmt"
)

type TomlDbRepository struct {
	driver TomlDriver
}

func NewTomlDbRepository(filepath string) (*TomlDbRepository, error) {
	driver := TomlDriver{
		filePath: filepath,
	}

	return &TomlDbRepository{
		driver: driver,
	}, nil
}

func (r *TomlDbRepository) Upsert(groupName entity.GroupName, key entity.Key, newValue entity.Value) error {
	store, err := r.getStore()
	if err != nil {
		return err
	}

	if _, exists := store[groupName]; !exists {
		return fmt.Errorf("group %s does not exists", groupName)
	}
	store[groupName][key] = newValue
	return r.save(store)
}

func (r *TomlDbRepository) List(groupName entity.GroupName) (entity.DbMap, error) {
	return r.getEntries(groupName)
}

func (r *TomlDbRepository) Delete(groupName entity.GroupName, key entity.Key) error {
	store, err := r.getStore()
	if err != nil {
		return err
	}

	delete(store[groupName], key)
	return r.save(store)
}

func (r *TomlDbRepository) getEntries(groupName entity.GroupName) (entity.DbMap, error) {
	store, err := r.getStore()
	if err != nil {
		return nil, err
	}
	entries, exists := store[groupName]
	if !exists {
		return nil, errors.New("could not find group named " + string(groupName))
	}
	return entries, nil
}

func (r *TomlDbRepository) getStore() (entity.DbSnapshot, error) {
	storage, err := r.driver.Load()
	if err != nil {
		return nil, err
	}
	return map_storage_to_entity(storage), nil
}

func (r *TomlDbRepository) save(storeGroup entity.DbSnapshot) error {
	storage := map_entity_to_storage(storeGroup)
	return r.driver.Save(storage)
}
