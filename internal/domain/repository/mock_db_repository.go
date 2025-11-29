package repository

import "duh/internal/domain/entity"

type MockInmemoryDbRepository struct {
	snapshot entity.DbSnapshot
}

func NewMockInmemoryDbRepository() *MockInmemoryDbRepository {
	return &MockInmemoryDbRepository{
		snapshot: make(entity.DbSnapshot),
	}
}

func (repo *MockInmemoryDbRepository) Upsert(groupName entity.GroupName, key entity.Key, value entity.Value) error {
	if _, exists := repo.snapshot[groupName]; !exists {
		repo.snapshot[groupName] = make(entity.DbMap)
	}
	repo.snapshot[groupName][key] = value
	return nil
}

func (repo *MockInmemoryDbRepository) List(groupName entity.GroupName) (entity.DbMap, error) {
	if entries, exists := repo.snapshot[groupName]; exists {
		return entries, nil
	}
	return make(entity.DbMap), nil
}

func (repo *MockInmemoryDbRepository) Delete(groupName entity.GroupName, key entity.Key) error {
	if entries, exists := repo.snapshot[groupName]; exists {
		delete(entries, key)
	}
	return nil
}
