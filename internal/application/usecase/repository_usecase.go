package usecase

import (
	"duh/internal/domain/entity"
	"duh/internal/domain/service"
)

type RepositoryUsecase struct {
	repositoryService *service.RepositoryService
}

func NewRepositoryUsecase(repositoryService *service.RepositoryService) *RepositoryUsecase {
	return &RepositoryUsecase{
		repositoryService: repositoryService,
	}
}

func (r *RepositoryUsecase) ListRepositories() (map[string][]string, error) {
	// Delegate to domain service for business logic
	return r.repositoryService.GetRepositoriesGroupedByStatus()
}

func (r *RepositoryUsecase) EnableRepository(repoName string) error {
	// Delegate to domain service for business logic
	return r.repositoryService.EnableRepository(repoName)
}

func (r *RepositoryUsecase) DisableRepository(repoName string) error {
	// Delegate to domain service for business logic
	return r.repositoryService.DisableRepository(repoName)
}

func (r *RepositoryUsecase) DeleteRepository(repoName string) error {
	// Application layer: orchestrate disable then delete
	if err := r.repositoryService.DisableRepository(repoName); err != nil {
		// If disable fails due to business rules, still allow deletion
	}
	return r.repositoryService.DeleteRepository(repoName)
}

func (r *RepositoryUsecase) SetDefaultRepository(repoName string) error {
	// Application layer: orchestrate enable then set default
	return r.repositoryService.SetDefaultRepository(repoName)
}

func (r *RepositoryUsecase) GetDefaultRepository() (string, error) {
	// Simple delegation
	return r.repositoryService.GetDefaultRepositoryName()
}

func (r *RepositoryUsecase) RenameRepository(oldName, newName string) error {
	// Delegate to domain service
	return r.repositoryService.RenameRepository(oldName, newName)
}

func (r *RepositoryUsecase) AddRepository(url string, name *string) error {
	// Application layer: orchestrate add then enable
	return r.repositoryService.AddAndEnableRepository(url, name)
}

func (r *RepositoryUsecase) CreateRepository(name string) error {
	// Delegate to domain service
	return r.repositoryService.CreateRepository(name)
}

func (r *RepositoryUsecase) UpdateRepositories(strategy string) (entity.RepositoryUpdateResults, error) {
	// Delegate to domain service
	return r.repositoryService.UpdateRepositories(strategy)
}

func (r *RepositoryUsecase) EditRepository(repoName string) error {
	// Delegate to domain service
	return r.repositoryService.EditRepository(repoName)
}

func (r *RepositoryUsecase) PushRepository(repoName string) error {
	// Delegate to domain service
	return r.repositoryService.PushRepository(repoName)
}

func (r *RepositoryUsecase) EditGitconfig(repoName string) error {
	// Delegate to domain service
	return r.repositoryService.EditGitconfig(repoName)
}
