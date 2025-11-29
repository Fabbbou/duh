package service

import (
	"duh/internal/domain/repository"
	"slices"
)

type RepositoriesService struct {
	directoryService    *DirectoryService
	dbRepositoryFactory repository.DbRepositoryFactory
	allRepositories     map[string]repository.DbRepository
}

func NewRepositoriesService(directoryService *DirectoryService, dbRepositoryFactory repository.DbRepositoryFactory) *RepositoriesService {
	return &RepositoriesService{
		directoryService:    directoryService,
		dbRepositoryFactory: dbRepositoryFactory,
		allRepositories:     make(map[string]repository.DbRepository),
	}
}

func (svc *RepositoriesService) LoadExistingRepositories() error {
	repoPaths, err := svc.directoryService.ListRepositoryNames()
	if err != nil {
		return err
	}
	for _, repoPath := range repoPaths {
		dbRepo, err := svc.dbRepositoryFactory.NewDbRepository(repoPath)
		if err != nil {
			return err
		}
		svc.allRepositories[repoPath] = dbRepo
	}
	return nil
}

func (svc *RepositoriesService) GetRepositories(filter []string) map[string]repository.DbRepository {
	repos := make(map[string]repository.DbRepository)
	if len(filter) == 0 {
		return repos
	}
	for name, repo := range svc.allRepositories {
		if slices.Contains(filter, name) {
			repos[name] = repo
		}
	}
	return repos
}
