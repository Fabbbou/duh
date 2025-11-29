package toml_storage

import (
	"duh/internal/domain/entity"
	"errors"
	"fmt"
)

type TomlStoreRepository struct {
	driver TomlDriver
}

func NewTomlStoreRepository(filepath string) (*TomlStoreRepository, error) {
	driver := TomlDriver{
		filePath: filepath,
	}

	return &TomlStoreRepository{
		driver: driver,
	}, nil
}

func (r *TomlStoreRepository) Upsert(groupName entity.GroupName, key entity.Key, newValue entity.Value) error {
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

func (r *TomlStoreRepository) List(groupName entity.GroupName) (entity.StoreEntries, error) {
	return r.getEntries(groupName)
}

func (r *TomlStoreRepository) Delete(groupName entity.GroupName, key entity.Key) error {
	store, err := r.getStore()
	if err != nil {
		return err
	}

	delete(store[groupName], key)
	return r.save(store)
}

func (r *TomlStoreRepository) getEntries(groupName entity.GroupName) (entity.StoreEntries, error) {
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

func (r *TomlStoreRepository) getStore() (entity.Store, error) {
	storage, err := r.driver.Load()
	if err != nil {
		return nil, err
	}
	return map_storage_to_entity(storage), nil
}

func (r *TomlStoreRepository) save(storeGroup entity.Store) error {
	storage := map_entity_to_storage(storeGroup)
	return r.driver.Save(storage)
}
