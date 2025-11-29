package service

import (
	"duh/internal/domain/entity"
	"duh/internal/domain/repository"
	"duh/internal/domain/utils"
	"fmt"
	"strings"
)

type CliService struct {
	reposService          *RepositoriesService
	userPreferenceService UserPreferenceService
}

func NewCliService(
	reposService *RepositoriesService,
	userPreferenceService UserPreferenceService,
) CliService {
	return CliService{
		reposService:          reposService,
		userPreferenceService: userPreferenceService,
	}
}

func (svc *CliService) Inject() (string, error) {
	activatedReposNames, err := svc.userPreferenceService.GetCurrentActivatedRepositories()
	activatedRepos := svc.reposService.GetRepositories(activatedReposNames)
	if err != nil {
		return "", err
	}

	injectionLines := make([]string, 0)
	for _, repo := range activatedRepos {
		aliases, err := repo.List(entity.Aliases)
		if err != nil {
			return "", err
		}
		for key, value := range aliases {
			line, err := svc.buildInjectionLine(entity.Aliases, key, value)
			if err != nil {
				return "", err
			}
			injectionLines = append(injectionLines, line)
		}
		exports, err := repo.List(entity.Exports)
		if err != nil {
			return "", err
		}
		for key, value := range exports {
			line, err := svc.buildInjectionLine(entity.Exports, key, value)
			if err != nil {
				return "", err
			}
			injectionLines = append(injectionLines, line)
		}
	}
	lines := strings.Join(injectionLines, "\n")
	return lines, nil
}

func (svc *CliService) buildInjectionLine(group entity.GroupName, key entity.Key, value entity.Value) (string, error) {
	switch group {
	case entity.Aliases:
		return fmt.Sprintf(`alias %s="%s"`, key, utils.EnsureEscapeDoubleQuotes(string(value))), nil
	case entity.Exports:
		return fmt.Sprintf(`export %s="%s"`, key, utils.EnsureEscapeDoubleQuotes(string(value))), nil
	default:
		return "", fmt.Errorf("unknown group: %s", group)
	}
}

func (svc *CliService) getCurrentDefaultRepo() (repository.DbRepository, error) {
	defaultRepoName, err := svc.userPreferenceService.GetDefaultRepoName()
	if err != nil {
		return nil, err
	}
	repo, exists := svc.reposService.allRepositories[defaultRepoName]
	if !exists {
		return nil, fmt.Errorf("default repository %s not found", defaultRepoName)
	}
	return repo, nil
}

func (svc *CliService) AddAlias(key entity.Key, value entity.Value) error {
	repo, err := svc.getCurrentDefaultRepo()
	if err != nil {
		return err
	}
	return repo.Upsert(entity.Aliases, key, value)
}

func (svc *CliService) RemoveAlias(key entity.Key) error {
	repo, err := svc.getCurrentDefaultRepo()
	if err != nil {
		return err
	}
	return repo.Delete(entity.Aliases, key)
}

func (svc *CliService) ListAliases() (entity.DbMap, error) {
	repo, err := svc.getCurrentDefaultRepo()
	if err != nil {
		return nil, err
	}
	return repo.List(entity.Aliases)
}

func (svc *CliService) AddExport(key entity.Key, value entity.Value) error {
	repo, err := svc.getCurrentDefaultRepo()
	if err != nil {
		return err
	}
	return repo.Upsert(entity.Exports, key, value)
}

func (svc *CliService) RemoveExport(key entity.Key) error {
	repo, err := svc.getCurrentDefaultRepo()
	if err != nil {
		return err
	}
	return repo.Delete(entity.Exports, key)
}

func (svc *CliService) ListExports() (entity.DbMap, error) {
	repo, err := svc.getCurrentDefaultRepo()
	if err != nil {
		return nil, err
	}
	return repo.List(entity.Exports)
}
