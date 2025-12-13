package file_repository

import "duh/internal/domain/entity"

type RepositoryRepository interface {
	Save(newVersion entity.Repository) error
	GetDbPath() string
	Get() (entity.Repository, error)
}
