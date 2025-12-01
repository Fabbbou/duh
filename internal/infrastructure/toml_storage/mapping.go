package toml_storage

import (
	"duh/internal/domain/entity"
	"fmt"
)

func map_to_entity(payload interface{}) (entity.DbSnapshot, error) {
	if storage, ok := payload.(*RepositoryDb); ok {
		return map_repodb_to_entity(storage), nil
	} else if userpref, ok := payload.(*UserPreferenceDb); ok {
		return map_userpref_to_entity(*userpref), nil
	}
	return entity.DbSnapshot{}, fmt.Errorf("could not find the right db type to parse for <%s>", payload)
}

func map_from_entity(groupName entity.GroupName, store entity.DbSnapshot) (interface{}, error) {
	switch groupName {
	case entity.Repositories:
		return map_entity_to_userprefdb(store), nil
	case entity.Aliases, entity.Exports:
		return map_entity_to_repodb(store)
	}
	return nil, fmt.Errorf("could not find the right db type to parse for group <%s>", groupName)
}

func map_entity_to_userprefdb(store entity.DbSnapshot) *UserPreferenceDb {
	repositories, exists := store[entity.Repositories]
	if !exists {
		repositories = entity.DbMap{}
	}

	userprefdb := &UserPreferenceDb{
		Repositories: make(map[string]string, len(repositories)),
	}

	for key, value := range repositories {
		userprefdb.Repositories[key] = value
	}

	return userprefdb
}

func map_userpref_to_entity(payload UserPreferenceDb) entity.DbSnapshot {
	result := make(entity.DbSnapshot)

	entries := make(entity.DbMap, len(payload.Repositories))
	for k, v := range payload.Repositories {
		entries[k] = v
	}
	result[entity.Repositories] = entries
	return result
}

func map_repodb_to_entity(storage *RepositoryDb) entity.DbSnapshot {
	result := make(entity.DbSnapshot)

	entries := make(entity.DbMap, len(storage.Aliases))
	for k, v := range storage.Aliases {
		entries[k] = v
	}
	result[entity.Aliases] = entries

	entries = make(entity.DbMap, len(storage.Exports))
	for k, v := range storage.Exports {
		entries[k] = v
	}
	result[entity.Exports] = entries

	return result
}

func map_entity_to_repodb(store entity.DbSnapshot) (*RepositoryDb, error) {
	for k := range store {
		if k != entity.Aliases && k != entity.Exports {
			return nil, fmt.Errorf("cannot map group <%s> to RepositoryDb", k)
		}
	}

	aliases, exists := store[entity.Aliases]
	if !exists {
		aliases = entity.DbMap{}
	}

	exports, exists := store[entity.Exports]
	if !exists {
		exports = entity.DbMap{}
	}

	storage := &RepositoryDb{
		Aliases: make(map[string]string, len(aliases)),
		Exports: make(map[string]string, len(exports)),
	}

	for key, value := range aliases {
		storage.Aliases[key] = value
	}

	for key, value := range exports {
		storage.Exports[key] = value
	}

	return storage, nil
}
