package repository

import "duh/internal/domain/entity"

type DbRepository interface {

	/// Get all enabled repositories
	GetEnabledRepositories() ([]entity.Repository, error)

	/// Get the default repository
	GetDefaultRepository() (*entity.Repository, error)

	/// Add or update a repository
	UpsertRepository(repo entity.Repository) error

	/// List all repositories
	GetAllRepositories() ([]entity.Repository, error)

	/// Delete a repository
	DeleteRepository(repoName string) error

	/// Rename a repository
	// RenameRepository(oldName, newName string) error

	/// Set a repository as the default one
	ChangeDefaultRepository(repoName string) error

	/// Disable a repository from being used
	DisableRepository(repoName string) error

	/// Enable a repository to be used
	EnableRepository(repoName string) error

	/// Initialiaze the database if needed
	CheckInit() (bool, error)
}
