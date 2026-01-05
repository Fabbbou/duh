package usecase

import (
	"duh/internal/domain/entity"
	"duh/internal/domain/port"
)

type RepositoryUsecase struct {
	dbPort port.DbPort
}

func NewRepositoryUsecase(dbPort port.DbPort) *RepositoryUsecase {
	return &RepositoryUsecase{
		dbPort: dbPort,
	}
}

func (r *RepositoryUsecase) ListRepositories() (map[string][]string, error) {
	repos, err := r.dbPort.GetAllRepositories()
	if err != nil {
		return nil, err
	}

	enabledRepos, err := r.dbPort.GetEnabledRepositories()
	if err != nil {
		return nil, err
	}

	enabledMap := make(map[string]bool)
	for _, repo := range enabledRepos {
		enabledMap[repo.Name] = true
	}

	result := map[string][]string{
		"enabled":  {},
		"disabled": {},
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

func (r *RepositoryUsecase) EnableRepository(repoName string) error {
	return r.dbPort.EnableRepository(repoName)
}

func (r *RepositoryUsecase) DisableRepository(repoName string) error {
	return r.dbPort.DisableRepository(repoName)
}

func (r *RepositoryUsecase) DeleteRepository(repoName string) error {
	return r.dbPort.DeleteRepository(repoName)
}

func (r *RepositoryUsecase) SetDefaultRepository(repoName string) error {
	err := r.dbPort.EnableRepository(repoName)
	if err != nil {
		return err
	}
	return r.dbPort.ChangeDefaultRepository(repoName)
}

func (r *RepositoryUsecase) GetDefaultRepository() (string, error) {
	repo, err := r.dbPort.GetDefaultRepository()
	if err != nil {
		return "", err
	}
	return repo.Name, nil
}

func (r *RepositoryUsecase) RenameRepository(oldName, newName string) error {
	return r.dbPort.RenameRepository(oldName, newName)
}

func (r *RepositoryUsecase) AddRepository(url string, name *string) error {
	repo, err := r.dbPort.AddRepository(url, name)
	if err != nil {
		return err
	}
	return r.dbPort.EnableRepository(repo)
}

func (r *RepositoryUsecase) CreateRepository(name string) error {
	_, err := r.dbPort.CreateRepository(name)
	return err
}

func (r *RepositoryUsecase) UpdateRepositories(strategy string) (entity.RepositoryUpdateResults, error) {
	return r.dbPort.UpdateRepositories(strategy)
}

func (r *RepositoryUsecase) EditRepository(repoName string) error {
	return r.dbPort.EditRepo(repoName)
}

func (r *RepositoryUsecase) PushRepository(repoName string) error {
	return r.dbPort.PushRepository(repoName)
}

func (r *RepositoryUsecase) EditGitconfig(repoName string) error {
	return r.dbPort.EditGitconfig(repoName)
}
