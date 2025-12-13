package file_db

import (
	"duh/internal/domain/entity"
	"fmt"
)

type TomlRepositoryRepository struct {
	tomlDriver *TomlDriver
}

func NewTomlRepositoryRepository(dbFilePath string) *TomlRepositoryRepository {
	return &TomlRepositoryRepository{
		tomlDriver: &TomlDriver{dbFilePath},
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
	repositoryDb, ok := rawData.(*RepositoryDb)
	if !ok {
		return entity.Repository{}, fmt.Errorf("could not cast repository file from data in file %s", r.GetDbPath())
	}
	return entity.Repository{
		Aliases: repositoryDb.Aliases,
		Exports: repositoryDb.Exports,
	}, nil
}
