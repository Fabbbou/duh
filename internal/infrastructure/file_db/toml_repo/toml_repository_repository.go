package toml_repo

import (
	"duh/internal/domain/entity"
)

type TomlRepositoryRepository struct {
	tomlDriver *TomlDriver[RepositoryDb]
}

// could just be a static function that returns a repository entity instead
func NewTomlRepositoryRepository(dbFilePath string) *TomlRepositoryRepository {
	return &TomlRepositoryRepository{
		tomlDriver: &TomlDriver[RepositoryDb]{dbFilePath},
	}
}

func (r *TomlRepositoryRepository) Save(newVersion entity.Repository) error {
	//map to toml compatible struct
	repositoryDb := RepositoryDb{
		Aliases: newVersion.Aliases,
		Exports: newVersion.Exports,
	}
	return r.tomlDriver.Save(repositoryDb)
}

func (r *TomlRepositoryRepository) GetDbPath() string {
	return r.tomlDriver.filePath
}

func (r *TomlRepositoryRepository) Get() (entity.Repository, error) {
	rawData, err := r.tomlDriver.Load()
	if err != nil {
		return entity.Repository{}, err
	}
	repositoryDb := *rawData
	return entity.Repository{
		Aliases: repositoryDb.Aliases,
		Exports: repositoryDb.Exports,
	}, nil
}
