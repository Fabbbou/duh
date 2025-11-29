package service

import (
	"duh/internal/domain/entity"
	"duh/internal/domain/repository"
)

type GroupStoreService struct {
	groupName       entity.GroupName
	storeRepository repository.GroupStoreRepository
}

func NewGroupStoreService(groupName entity.GroupName, groupRepository repository.GroupStoreRepository) *GroupStoreService {
	return &GroupStoreService{
		groupName:       groupName,
		storeRepository: groupRepository,
	}
}

func (gs *GroupStoreService) List() (entity.StoreEntries, error) {
	return gs.storeRepository.List(gs.groupName)
}

func (gs *GroupStoreService) Upsert(key entity.Key, value entity.Value) error {
	return gs.storeRepository.Upsert(gs.groupName, key, value)
}

func (gs *GroupStoreService) Delete(key entity.Key) error {
	return gs.storeRepository.Delete(gs.groupName, key)
}
