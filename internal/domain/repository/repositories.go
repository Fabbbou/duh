package repository

import "duh/internal/domain/entity"

type GroupStoreRepository interface {
	Upsert(groupName entity.GroupName, key entity.Key, value entity.Value) error
	List(groupName entity.GroupName) (entity.StoreEntries, error)
	Delete(groupName entity.GroupName, key entity.Key) error
}

type CliQueryRepository interface {
	/// Inject all the store entries into the CLI context
	/// eval "$(duh inject)"
	/// or
	/// source <(duh inject)
	/// or (if duh is already installed)
	/// duh_reload
	Inject(entity.Store) error

	// AddAlias(key entity.Key, value entity.Value) error
	// RemoveAlias(key entity.Key) error
	// ListAliases() (entity.StoreEntries, error)

	// AddExport(key entity.Key, value entity.Value) error
	// RemoveExport(key entity.Key) error
	// ListExports() (entity.StoreEntries, error)
}
