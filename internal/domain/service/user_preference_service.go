package service

import (
	"duh/internal/domain/entity"
	"duh/internal/domain/repository"
	"duh/internal/domain/utils"
	"errors"
)

type UserPreferenceService struct {
	userDbRepo repository.DbRepository
}

const RepositoriesGroup entity.GroupName = "repositories"

const ActivatedRepositoriesKey entity.Key = "activated"
const LocalRepoNameKey entity.Key = "local_repo_name"
const DefaultRepoName entity.Key = "default"

func NewUserPreferenceService(
	userDbRepo repository.DbRepository,
) UserPreferenceService {
	return UserPreferenceService{
		userDbRepo: userDbRepo,
	}
}

func (svc *UserPreferenceService) GetCurrentActivatedRepositories() ([]string, error) {
	repositories, err := svc.userDbRepo.List(RepositoriesGroup)
	if err != nil {
		return nil, err
	}
	activatedReposValue, exists := repositories[ActivatedRepositoriesKey]
	if !exists || activatedReposValue == "" {
		return []string{}, nil
	}
	return utils.ParseCommaSeparatedValues(activatedReposValue), nil
}

func (svc *UserPreferenceService) AddActivatedRepository(repoName string) error {
	currentRepos, err := svc.GetCurrentActivatedRepositories()
	if err != nil {
		return err
	}
	for _, r := range currentRepos {
		if r == repoName {
			// Already activated
			return nil
		}
	}
	currentRepos = append(currentRepos, repoName)
	updatedValue := entity.Value(utils.EnsureEscapeDoubleQuotes(
		utils.JoinCommaSeparatedValues(currentRepos),
	))
	return svc.userDbRepo.Upsert(
		RepositoriesGroup,
		ActivatedRepositoriesKey,
		updatedValue,
	)
}

func (svc *UserPreferenceService) RemoveActivatedRepository(repoName string) error {
	currentRepos, err := svc.GetCurrentActivatedRepositories()
	if err != nil {
		return err
	}
	updatedRepos := []string{}
	for _, r := range currentRepos {
		if r != repoName {
			updatedRepos = append(updatedRepos, r)
		}
	}
	updatedValue := entity.Value(utils.EnsureEscapeDoubleQuotes(
		utils.JoinCommaSeparatedValues(updatedRepos),
	))
	return svc.userDbRepo.Upsert(
		RepositoriesGroup,
		ActivatedRepositoriesKey,
		updatedValue,
	)
}

func (svc *UserPreferenceService) setKeyRepo(key string, value string) error {
	return svc.userDbRepo.Upsert(
		RepositoriesGroup,
		key,
		entity.Value(utils.EnsureEscapeDoubleQuotes(value)),
	)
}

func (svc *UserPreferenceService) getKeyRepo(key string) (string, error) {
	repositories, err := svc.userDbRepo.List(RepositoriesGroup)
	if err != nil {
		return "", err
	}
	localRepoNameValue, exists := repositories[key]
	if !exists {
		return "", nil
	}
	return string(localRepoNameValue), nil
}

func (svc *UserPreferenceService) SetLocalRepoName(name string) error {
	return svc.setKeyRepo(LocalRepoNameKey, name)
}

func (svc *UserPreferenceService) GetLocalRepoName() (string, error) {
	return svc.getKeyRepo(LocalRepoNameKey)
}

func (svc *UserPreferenceService) SetDefaultRepoName(name string) error {
	return svc.setKeyRepo(DefaultRepoName, name)
}

func (svc *UserPreferenceService) GetDefaultRepoName() (string, error) {
	return svc.getKeyRepo(DefaultRepoName)
}

func (svc *UserPreferenceService) InitUserPreference() error {

	repositories, err := svc.userDbRepo.List(RepositoriesGroup)
	if err != nil {
		return err
	}
	_, exists := repositories[LocalRepoNameKey]
	if exists {
		return errors.New("user preference already initialized")
	}
	err = svc.userDbRepo.Upsert(RepositoriesGroup, LocalRepoNameKey, "local")
	if err != nil {
		return err
	}
	err = svc.userDbRepo.Upsert(RepositoriesGroup, ActivatedRepositoriesKey, "local")
	if err != nil {
		return err
	}
	return svc.userDbRepo.Upsert(RepositoriesGroup, DefaultRepoName, "local")
}
