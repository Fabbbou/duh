package file_repository

import (
	"duh/internal/infrastructure/file_db/file_dto"
)

type UserPreferencesRepository interface {
	Save(file_dto.UserPreferences) error
	GetDbPath() string
	Get() (file_dto.UserPreferences, error)
}
