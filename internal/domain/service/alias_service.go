package service

import (
	"duh/internal/domain/entity"
	"duh/internal/domain/port"
	"maps"
)

type AliasService struct {
	dbPort port.DbPort
}

func NewAliasService(dbPort port.DbPort) *AliasService {
	return &AliasService{
		dbPort: dbPort,
	}
}

// SetAlias sets an alias in the default repository
func (a *AliasService) SetAlias(aliasName, value string) error {
	repo, err := a.dbPort.GetDefaultRepository()
	if err != nil {
		return err
	}

	// Business logic: initialize aliases map if nil
	if repo.Aliases == nil {
		repo.Aliases = make(map[string]string)
	}

	repo.Aliases[aliasName] = value
	return a.dbPort.UpsertRepository(*repo)
}

// UnsetAlias removes an alias from the default repository
func (a *AliasService) UnsetAlias(aliasName string) error {
	repo, err := a.dbPort.GetDefaultRepository()
	if err != nil {
		return err
	}

	// Business logic: only delete if aliases map exists
	if repo.Aliases != nil {
		delete(repo.Aliases, aliasName)
	}

	return a.dbPort.UpsertRepository(*repo)
}

// GetMergedAliases aggregates aliases from all enabled repositories
func (a *AliasService) GetMergedAliases() (map[string]string, error) {
	repos, err := a.dbPort.GetEnabledRepositories()
	if err != nil {
		return nil, err
	}

	// Business logic: merge aliases with priority (later repos override earlier ones)
	entries := make(map[string]string)
	for _, repo := range repos {
		if repo.Aliases != nil {
			maps.Copy(entries, repo.Aliases)
		}
	}

	return entries, nil
}

// ValidateAliasName checks if an alias name is valid according to business rules
func (a *AliasService) ValidateAliasName(aliasName string) error {
	if aliasName == "" {
		return &entity.ValidationError{Field: "alias_name", Message: "alias name cannot be empty"}
	}

	// Business rule: alias names should not contain spaces
	// Add more validation as needed

	return nil
}
