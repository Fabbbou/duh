package usecase

import (
	"duh/internal/domain/port"
	"maps"
)

type ExportsUsecase struct {
	dbPort port.DbPort
}

func NewExportsUsecase(dbPort port.DbPort) *ExportsUsecase {
	return &ExportsUsecase{
		dbPort: dbPort,
	}
}

func (e *ExportsUsecase) SetExport(exportName, value string) error {
	repo, err := e.dbPort.GetDefaultRepository()
	if err != nil {
		return err
	}
	repo.Exports[exportName] = value
	return e.dbPort.UpsertRepository(*repo)
}

func (e *ExportsUsecase) UnsetExport(exportName string) error {
	repo, err := e.dbPort.GetDefaultRepository()
	if err != nil {
		return err
	}

	delete(repo.Exports, exportName)
	return e.dbPort.UpsertRepository(*repo)
}

func (e *ExportsUsecase) ListExports() (map[string]string, error) {
	repos, err := e.dbPort.GetEnabledRepositories()
	if err != nil {
		return nil, err
	}

	allExports := make(map[string]string)
	for _, repo := range repos {
		maps.Copy(allExports, repo.Exports)
	}
	return allExports, nil
}
