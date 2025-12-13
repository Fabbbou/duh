package service

import (
	"duh/internal/domain/entity"
	"duh/internal/domain/repository"
)

type CliService struct {
	dbRepository repository.DbRepository
}

func NewCliService(dbRepository repository.DbRepository) CliService {
	return CliService{
		dbRepository: dbRepository,
	}
}

func (svc *CliService) Inject() (string, error) {

	return "", nil
}

func (svc *CliService) getActivatedRepositories() ([]entity.Repository, error) {
	// userPrefs, err := svc.userPreferencesRepository.Get()
	// if err != nil {
	// 	return nil, err
	// }

	// You might want to add logic here to convert userPrefs.ActivatedRepositories to []entity.Repository
	return nil, nil
}

func (svc *CliService) getCurrentDefaultRepo() (*entity.Repository, error) {
	return nil, nil
}

func (svc *CliService) SetAlias(key string, value string) error {
	return nil
}

func (svc *CliService) RemoveAlias(key string) error {
	return nil
}

func (svc *CliService) ListAliases() (*entity.Repository, error) {
	return nil, nil
}

func (svc *CliService) AddExport(key string, value string) error {
	return nil
}

func (svc *CliService) RemoveExport(key string) error {
	return nil
}

func (svc *CliService) ListExports() (*entity.Repository, error) {
	return nil, nil
}
