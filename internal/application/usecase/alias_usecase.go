package usecase

import (
	"duh/internal/domain/port"
	"maps"
)

type AliasUsecase struct {
	dbPort port.DbPort
}

func NewAliasUsecase(dbPort port.DbPort) *AliasUsecase {
	return &AliasUsecase{
		dbPort: dbPort,
	}
}

func (a *AliasUsecase) SetAlias(aliasName, value string) error {
	repo, err := a.dbPort.GetDefaultRepository()
	if err != nil {
		return err
	}
	repo.Aliases[aliasName] = value
	return a.dbPort.UpsertRepository(*repo)
}

func (a *AliasUsecase) UnsetAlias(aliasName string) error {
	repo, err := a.dbPort.GetDefaultRepository()
	if err != nil {
		return err
	}

	delete(repo.Aliases, aliasName)
	return a.dbPort.UpsertRepository(*repo)
}

func (a *AliasUsecase) ListAliases() (map[string]string, error) {
	repos, err := a.dbPort.GetEnabledRepositories()
	if err != nil {
		return nil, err
	}
	entries := map[string]string{}
	for _, repo := range repos {
		maps.Copy(entries, repo.Aliases)
	}
	return entries, nil
}
