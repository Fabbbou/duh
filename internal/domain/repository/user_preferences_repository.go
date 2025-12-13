package repository

import "duh/internal/domain/entity"

type UserPreferencesRepository interface {
	Save(entity.UserPreferences) error
	GetDbPath() string
	Get() (entity.UserPreferences, error)
}
