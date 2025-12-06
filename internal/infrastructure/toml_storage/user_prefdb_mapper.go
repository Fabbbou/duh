package toml_storage

import (
	"duh/internal/domain/entity"
	"fmt"
)

type MapperType string

type EntityMapper interface {
	ToEntity(payload interface{}) (*entity.DbSnapshot, error)
	ToDb(entity entity.DbSnapshot) (interface{}, error)
}

type UserPrefDbMapper struct{}

func NewUserPrefDbMapper() *UserPrefDbMapper {
	return &UserPrefDbMapper{}
}

func (m *UserPrefDbMapper) ToEntity(payload interface{}) (*entity.DbSnapshot, error) {
	payloadTyped, ok := payload.(UserPreferenceDb)
	if !ok {
		return nil, fmt.Errorf("invalid payload type for UserPreferenceDbMapper: %T", payload)
	}

	result := make(entity.DbSnapshot)

	entries := make(entity.DbMap, len(payloadTyped.Repositories))
	for k, v := range payloadTyped.Repositories {
		entries[k] = v
	}
	result[entity.Repositories] = entries
	return &result, nil
}

func (m *UserPrefDbMapper) ToDb(entit entity.DbSnapshot) (interface{}, error) {
	repositories, exists := entit[entity.Repositories]
	if !exists {
		repositories = entity.DbMap{}
	}

	userprefdb := &UserPreferenceDb{
		Repositories: make(map[string]string, len(repositories)),
	}

	for key, value := range repositories {
		userprefdb.Repositories[key] = value
	}

	return *userprefdb, nil
}
