package toml_storage

import (
	"duh/internal/domain/entity"
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
		store[groupName] = make(entity.DbMap)
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
		return entity.DbMap{}, nil
	}
	return entries, nil
}

func (r *TomlDbRepository) getStore() (entity.DbSnapshot, error) {
	storage, err := r.driver.Load()
	if err != nil {
		return nil, err
	}
	return map_to_entity(storage)
}

func (r *TomlDbRepository) save(storeGroup entity.DbSnapshot) error {
	if _, exists := storeGroup[entity.Repositories]; exists {
		storage, err := map_from_entity(entity.Repositories, storeGroup)
		if err != nil {
			return err
		}
		return r.driver.Save(storage)
	}
	_, existsAliases := storeGroup[entity.Aliases]
	_, existsExports := storeGroup[entity.Exports]
	if existsExports || existsAliases {
		storage, err := map_from_entity(entity.Aliases, storeGroup)
		if err != nil {
			return err
		}
		return r.driver.Save(storage)
	}
	return fmt.Errorf("could not find a correct format to save")
}
