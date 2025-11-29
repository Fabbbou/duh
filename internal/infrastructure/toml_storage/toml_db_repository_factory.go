package toml_storage

import "duh/internal/domain/repository"

type TomlDbRepositoryFactory struct{}

func NewTomlDbRepositoryFactory() *TomlDbRepositoryFactory {
	return &TomlDbRepositoryFactory{}
}

func (f *TomlDbRepositoryFactory) NewDbRepository(repositoryPath string) (repository.DbRepository, error) {
	return NewTomlDbRepository(repositoryPath)
}
