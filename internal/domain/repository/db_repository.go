package repository

import "duh/internal/domain/entity"

type DbRepository interface {
	SaveUserPreferences(prefs entity.UserPreferences) error
	GetUserPreferences() (entity.UserPreferences, error)
	SaveRepository(repoName string, repo entity.Repository) error
	GetActivatedRepositories() ([]entity.Repository, error)
	GetDefaultRepository() (*entity.Repository, error)
}
