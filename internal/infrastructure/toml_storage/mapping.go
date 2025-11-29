package toml_storage

import "duh/internal/domain/entity"

func map_storage_to_entity(storage *Storage) entity.Store {
	result := make(entity.Store)

	entries := make(entity.StoreEntries, len(storage.Aliases))
	for k, v := range storage.Aliases {
		entries[k] = v
	}
	result[entity.Aliases] = entries

	entries = make(entity.StoreEntries, len(storage.Exports))
	for k, v := range storage.Exports {
		entries[k] = v
	}
	result[entity.Exports] = entries

	return result
}

func map_entity_to_storage(store entity.Store) *Storage {
	aliases, exists := store[entity.Aliases]
	if !exists {
		aliases = entity.StoreEntries{}
	}

	exports, exists := store[entity.Exports]
	if !exists {
		exports = entity.StoreEntries{}
	}

	storage := &Storage{
		Aliases: make(map[string]string, len(aliases)),
		Exports: make(map[string]string, len(exports)),
	}

	for key, value := range aliases {
		storage.Aliases[key] = value
	}

	for key, value := range exports {
		storage.Exports[key] = value
	}

	return storage
}
