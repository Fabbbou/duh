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
}

type MockDbRepository struct {
	DefaultRepo entity.Repository
	Repos       []entity.Repository
	Enabled     []string
}

func (m *MockDbRepository) GetEnabledRepositories() ([]entity.Repository, error) {
	enabledRepos := []entity.Repository{}
	for _, repo := range m.Repos {
		for _, enabledName := range m.Enabled {
			if repo.Name == enabledName {
				enabledRepos = append(enabledRepos, repo)
			}
		}
	}
	return enabledRepos, nil
}

func (m *MockDbRepository) GetDefaultRepository() (*entity.Repository, error) {
	return &m.DefaultRepo, nil
}

func (m *MockDbRepository) UpsertRepository(repo entity.Repository) error {
	for i, r := range m.Repos {
		if r.Name == repo.Name {
			m.Repos[i] = repo
			return nil
		}
	}
	m.Repos = append(m.Repos, repo)
	return nil
}

func (m *MockDbRepository) GetAllRepositories() ([]entity.Repository, error) {
	return m.Repos, nil
}

func (m *MockDbRepository) DeleteRepository(repoName string) error {
	newRepos := []entity.Repository{}
	for _, r := range m.Repos {
		if r.Name != repoName {
			newRepos = append(newRepos, r)
		}
	}
	m.Repos = newRepos
	return nil
}

func (m *MockDbRepository) ChangeDefaultRepository(repoName string) error {
	for _, r := range m.Repos {
		if r.Name == repoName {
			m.DefaultRepo = r
			return nil
		}
	}
	return nil
}

func (m *MockDbRepository) DisableRepository(repoName string) error {
	newEnabled := []string{}
	for _, name := range m.Enabled {
		if name != repoName {
			newEnabled = append(newEnabled, name)
		}
	}
	m.Enabled = newEnabled
	return nil
}

func (m *MockDbRepository) EnableRepository(repoName string) error {
	for _, name := range m.Enabled {
		if name == repoName {
			return nil
		}
	}
	m.Enabled = append(m.Enabled, repoName)
	return nil
}

func (m *MockDbRepository) CheckInit() (bool, error) {
	return true, nil
}
