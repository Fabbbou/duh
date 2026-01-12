package port

import (
	"duh/internal/domain/entity"
	"fmt"
)

type DbPort interface {

	/// Get all enabled repositories
	GetEnabledPackages() ([]entity.Package, error)

	/// Get the default package
	GetDefaultPackage() (*entity.Package, error)

	/// Add or update a package
	UpsertPackage(repo entity.Package) error

	/// List all repositories
	GetAllPackages() ([]entity.Package, error)

	/// Delete a package
	DeletePackage(repoName string) error

	/// Rename a package
	RenamePackage(oldName, newName string) error

	/// Set a package as the default one
	ChangeDefaultPackage(repoName string) error

	/// Disable a package from being used
	DisablePackage(repoName string) error

	/// Enable a package to be used
	EnablePackage(repoName string) error

	/// Add a new package, optionally with a specified name
	// the string returned is the name of the added package
	AddPackage(url string, name *string) (string, error)

	// Create a new package with the given name
	// By default it will be enabled
	// Also returns the path to the created package
	CreatePackage(name string) (string, error)

	// Update repositories according to the specified strategy
	// Strategies:
	// - entity.UpdateSafe: Do not pull if local changes exist, return ErrChangesExist if changes are present
	// - entity.UpdateKeep: Commit local changes before pulling
	// - entity.UpdateForce: Discard local changes and reset to remote state
	UpdatePackages(strategy string) (entity.PackageUpdateResults, error)

	// Edit a package's configuration file using the system's default editor
	EditPackage(repoName string) error

	EditGitconfig(repoName string) error

	// Push local changes in a package to its remote
	PushPackage(repoName string) error

	// Get the base path(s) where repositories are stored
	GetBasePath() (string, error)

	ListPackagePath() ([]string, error)

	// Adding other things to injection if not related to repositories
	// It is used to inject other dependencies like gitconfig file path
	BonusInjection(enabledPackages []entity.Package) (string, error)
}

type MockDbAdapter struct {
	DefaultRepo entity.Package
	Packages    []entity.Package
	Enabled     []string
}

func (m *MockDbAdapter) GetEnabledPackages() ([]entity.Package, error) {
	enabledPackages := []entity.Package{}
	for _, repo := range m.Packages {
		for _, enabledName := range m.Enabled {
			if repo.Name == enabledName {
				enabledPackages = append(enabledPackages, repo)
			}
		}
	}
	return enabledPackages, nil
}

func (m *MockDbAdapter) GetDefaultPackage() (*entity.Package, error) {
	return &m.DefaultRepo, nil
}

func (m *MockDbAdapter) UpsertPackage(repo entity.Package) error {
	for i, r := range m.Packages {
		if r.Name == repo.Name {
			m.Packages[i] = repo
			return nil
		}
	}
	m.Packages = append(m.Packages, repo)
	return nil
}

func (m *MockDbAdapter) GetAllPackages() ([]entity.Package, error) {
	return m.Packages, nil
}

func (m *MockDbAdapter) DeletePackage(repoName string) error {
	newPackages := []entity.Package{}
	for _, r := range m.Packages {
		if r.Name != repoName {
			newPackages = append(newPackages, r)
		}
	}
	m.Packages = newPackages
	return nil
}

func (m *MockDbAdapter) ChangeDefaultPackage(repoName string) error {
	for _, r := range m.Packages {
		if r.Name == repoName {
			m.DefaultRepo = r
			return nil
		}
	}
	return nil
}

func (m *MockDbAdapter) DisablePackage(repoName string) error {
	newEnabled := []string{}
	for _, name := range m.Enabled {
		if name != repoName {
			newEnabled = append(newEnabled, name)
		}
	}
	m.Enabled = newEnabled
	return nil
}

func (m *MockDbAdapter) EnablePackage(repoName string) error {
	for _, name := range m.Enabled {
		if name == repoName {
			return nil
		}
	}
	m.Enabled = append(m.Enabled, repoName)
	return nil
}

func (m *MockDbAdapter) RenamePackage(oldName, newName string) error {
	// Update package in list
	for i, repo := range m.Packages {
		if repo.Name == oldName {
			m.Packages[i].Name = newName
			break
		}
	}

	// Update enabled list
	for i, name := range m.Enabled {
		if name == oldName {
			m.Enabled[i] = newName
			break
		}
	}

	// Update default repo if it's the one being renamed
	if m.DefaultRepo.Name == oldName {
		m.DefaultRepo.Name = newName
	}

	return nil
}

func (m *MockDbAdapter) AddPackage(url string, name *string) (string, error) {
	if m.Packages == nil {
		m.Packages = []entity.Package{}
	}
	if m.Enabled == nil {
		m.Enabled = []string{}
	}
	if name == nil {
		generatedName := "repo" + fmt.Sprint(len(m.Packages)+1)
		name = &generatedName
	}
	m.Packages = append(m.Packages, entity.Package{Name: *name})
	m.Enabled = append(m.Enabled, *name)
	return "test/" + *name, nil
}

func (m *MockDbAdapter) CheckInit() (bool, error) {
	return true, nil
}

func (m *MockDbAdapter) CreatePackage(name string) (string, error) {
	if m.Packages == nil {
		m.Packages = []entity.Package{}
	}
	if m.Enabled == nil {
		m.Enabled = []string{}
	}
	if name == "" {
		generatedName := "repo" + fmt.Sprint(len(m.Packages)+1)
		name = generatedName
	}
	m.Packages = append(m.Packages, entity.Package{Name: name})
	m.Enabled = append(m.Enabled, name)
	return "test/" + name, nil
}

func (m *MockDbAdapter) UpdatePackages(strategy string) (entity.PackageUpdateResults, error) {
	// Mock implementation does nothing
	return entity.PackageUpdateResults{}, nil
}

func (m *MockDbAdapter) EditPackage(repoName string) error {
	// Mock implementation does nothing
	return nil
}

func (m *MockDbAdapter) PushPackage(repoName string) error {
	// Mock implementation does nothing
	return nil
}

func (m *MockDbAdapter) GetBasePath() (string, error) {
	return "/home/user/.local/share/duh", nil
}

func (m *MockDbAdapter) ListPackagePath() ([]string, error) {
	return []string{
		"/home/user/.local/share/duh/repositories/default",
		"/home/user/.local/share/duh/repositories",
	}, nil
}

func (m *MockDbAdapter) BonusInjection(enabledPackages []entity.Package) (string, error) {
	return "", nil
}

func (f *MockDbAdapter) EditGitconfig(repoName string) error {
	return nil
}
