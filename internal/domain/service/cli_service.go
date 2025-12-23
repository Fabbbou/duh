package service

import (
	"duh/internal/domain/entity"
	"duh/internal/domain/repository"
	"duh/internal/domain/utils"
	"fmt"
	"strings"
)

type CliService struct {
	dbRepository repository.DbRepository
}

func NewCliService(dbRepository repository.DbRepository) CliService {
	return CliService{
		dbRepository: dbRepository,
	}
}

func (cli *CliService) Inject() (string, error) {
	enabledRepos, err := cli.dbRepository.GetEnabledRepositories()
	if err != nil {
		return "", err
	}
	injectionLines := []string{}
	for _, repo := range enabledRepos {
		for key, value := range repo.Aliases {
			escapedKey := utils.EscapeDoubleQuotes(key)
			escapedValue := utils.EscapeDoubleQuotes(value)
			injectionLines = append(injectionLines, fmt.Sprintf("alias %s=\"%s\"", escapedKey, escapedValue))
		}
		for key, value := range repo.Exports {
			escapedKey := utils.EscapeDoubleQuotes(key)
			escapedValue := utils.EscapeDoubleQuotes(value)
			injectionLines = append(injectionLines, fmt.Sprintf("export %s=\"%s\"", escapedKey, escapedValue))
		}
	}
	injectionString := strings.Join(injectionLines, "\n")
	return injectionString, nil
}

func (cli *CliService) SetAlias(key string, value string) error {
	return nil
}

func (cli *CliService) RemoveAlias(key string) error {
	return nil
}

func (cli *CliService) ListAliases() (*entity.Repository, error) {
	return nil, nil
}

func (cli *CliService) AddExport(key string, value string) error {
	return nil
}

func (cli *CliService) RemoveExport(key string) error {
	return nil
}

func (cli *CliService) ListExports() (*entity.Repository, error) {
	return nil, nil
}
