package repository

import "cmd/cli/main.go/internal/domain/entity"

type GroupStorageRepository interface {
	Upsert(groupName entity.GroupName, key entity.Key, value entity.Value) error
	List(groupName entity.GroupName) ([]entity.StorageEntry, error)
	Delete(groupName entity.GroupName, key entity.Key) error
}
