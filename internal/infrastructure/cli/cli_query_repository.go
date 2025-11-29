package cli

import (
	"duh/internal/domain/entity"
	"fmt"
)

type CliQueryRepository struct {
	commandProcessor *CommandProcessor
}

func NewCliQueryRepository(commandProcessor *CommandProcessor) *CliQueryRepository {
	return &CliQueryRepository{
		commandProcessor: commandProcessor,
	}
}

func (repo *CliQueryRepository) Inject(store entity.Store) error {
	for name, value := range store[entity.Aliases] {
		repo.commandProcessor.Exec("alias", fmt.Sprintf("%s=%s", name, value))
	}
	for name, value := range store[entity.Exports] {
		repo.commandProcessor.Exec("export", fmt.Sprintf("%s=%s", name, value))
	}
	return nil
}

func (repo *CliQueryRepository) AddAlias(key string, value string) error {
	// Implementation of AddAlias method
	return nil
}

func (repo *CliQueryRepository) RemoveAlias(key string) error {
	// Implementation of RemoveAlias method
	return nil
}

func (repo *CliQueryRepository) ListAliases() (map[string]string, error) {
	// Implementation of ListAliases method
	return nil, nil
}

func (repo *CliQueryRepository) AddExport(key string, value string) error {
	// Implementation of AddExport method
	return nil
}

func (repo *CliQueryRepository) RemoveExport(key string) error {
	// Implementation of RemoveExport method
	return nil
}

func (repo *CliQueryRepository) ListExports() (map[string]string, error) {
	// Implementation of ListExports method
	return nil, nil
}
