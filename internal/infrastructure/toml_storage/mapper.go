package toml_storage

import "cmd/cli/main.go/internal/domain/entity"

func map_storage_to_entity(storage *Storage) entity.StoreGroup {
	result := make(entity.StoreGroup)

	entries := make([]entity.StorageEntry, 0, len(storage.Aliases))
	for k, v := range storage.Aliases {
		entries = append(entries, entity.StorageEntry{
			Key:   k,
			Value: v,
		})
	}
	result[entity.Aliases] = entries

	entries = make([]entity.StorageEntry, 0, len(storage.Exports))
	for k, v := range storage.Exports {
		entries = append(entries, entity.StorageEntry{
			Key:   k,
			Value: v,
		})
	}
	result[entity.Exports] = entries

	return result
}
