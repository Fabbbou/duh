package service

import (
	"duh/internal/domain/entity"
	"duh/internal/domain/port"
)

type RepositoryService struct {
	dbPort port.DbPort
}

func NewRepositoryService(dbPort port.DbPort) *RepositoryService {
	return &RepositoryService{
		dbPort: dbPort,
	}
}

// GetRepositoriesGroupedByStatus returns repositories grouped by enabled/disabled status
func (r *RepositoryService) GetRepositoriesGroupedByStatus() (map[string][]string, error) {
	repos, err := r.dbPort.GetAllRepositories()
	if err != nil {
		return nil, err
	}

	enabledRepos, err := r.dbPort.GetEnabledRepositories()
	if err != nil {
		return nil, err
	}

	// Business logic: create a lookup map for enabled repositories
	enabledMap := make(map[string]bool, len(enabledRepos))
	for _, repo := range enabledRepos {
		enabledMap[repo.Name] = true
	}

	// Business logic: categorize repositories
	result := map[string][]string{
		"enabled":  make([]string, 0),
		"disabled": make([]string, 0),
	}

	for _, repo := range repos {
		if enabledMap[repo.Name] {
			result["enabled"] = append(result["enabled"], repo.Name)
		} else {
			result["disabled"] = append(result["disabled"], repo.Name)
		}
	}

	return result, nil
}

// EnableRepository enables a repository and validates business rules
func (r *RepositoryService) EnableRepository(repoName string) error {
	// Business rule: validate repository exists
	if err := r.validateRepositoryExists(repoName); err != nil {
		return err
	}

	return r.dbPort.EnableRepository(repoName)
}

// DisableRepository disables a repository and validates business rules
func (r *RepositoryService) DisableRepository(repoName string) error {
	// Business rule: validate repository exists
	if err := r.validateRepositoryExists(repoName); err != nil {
		return err
	}

	// Business rule: cannot disable the last enabled repository
	enabled, err := r.dbPort.GetEnabledRepositories()
	if err != nil {
		return err
	}

	if len(enabled) <= 1 {
		return &entity.BusinessRuleError{
			Rule:    "minimum_one_repository",
			Message: "cannot disable the last enabled repository",
		}
	}

	return r.dbPort.DisableRepository(repoName)
}

// validateRepositoryExists checks if a repository exists
func (r *RepositoryService) validateRepositoryExists(repoName string) error {
	repos, err := r.dbPort.GetAllRepositories()
	if err != nil {
		return err
	}

	for _, repo := range repos {
		if repo.Name == repoName {
			return nil
		}
	}

	return &entity.NotFoundError{
		Resource: "repository",
		ID:       repoName,
	}
}

// Additional methods needed by the use case

func (r *RepositoryService) DeleteRepository(repoName string) error {
	if err := r.validateRepositoryExists(repoName); err != nil {
		return err
	}
	return r.dbPort.DeleteRepository(repoName)
}

func (r *RepositoryService) SetDefaultRepository(repoName string) error {
	if err := r.EnableRepository(repoName); err != nil {
		return err
	}
	return r.dbPort.ChangeDefaultRepository(repoName)
}

func (r *RepositoryService) GetDefaultRepositoryName() (string, error) {
	repo, err := r.dbPort.GetDefaultRepository()
	if err != nil {
		return "", err
	}
	return repo.Name, nil
}

func (r *RepositoryService) RenameRepository(oldName, newName string) error {
	if err := r.validateRepositoryExists(oldName); err != nil {
		return err
	}
	return r.dbPort.RenameRepository(oldName, newName)
}

func (r *RepositoryService) AddAndEnableRepository(url string, name *string) error {
	repo, err := r.dbPort.AddRepository(url, name)
	if err != nil {
		return err
	}
	return r.dbPort.EnableRepository(repo)
}

func (r *RepositoryService) CreateRepository(name string) error {
	_, err := r.dbPort.CreateRepository(name)
	return err
}

func (r *RepositoryService) UpdateRepositories(strategy string) (entity.RepositoryUpdateResults, error) {
	return r.dbPort.UpdateRepositories(strategy)
}

func (r *RepositoryService) EditRepository(repoName string) error {
	if err := r.validateRepositoryExists(repoName); err != nil {
		return err
	}
	return r.dbPort.EditRepo(repoName)
}

func (r *RepositoryService) PushRepository(repoName string) error {
	if err := r.validateRepositoryExists(repoName); err != nil {
		return err
	}
	return r.dbPort.PushRepository(repoName)
}

func (r *RepositoryService) EditGitconfig(repoName string) error {
	if err := r.validateRepositoryExists(repoName); err != nil {
		return err
	}
	return r.dbPort.EditGitconfig(repoName)
}
