package repository

import "duh/internal/domain/entity"

type DbRepository interface {
	Upsert(groupName entity.GroupName, key entity.Key, value entity.Value) error
	List(groupName entity.GroupName) (entity.DbMap, error)
	Delete(groupName entity.GroupName, key entity.Key) error
}

type DbRepositoryFactory interface {
	NewDbRepository(repositoryPath string) (DbRepository, error)
}
