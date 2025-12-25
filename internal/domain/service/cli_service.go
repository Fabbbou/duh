package service

import (
	"duh/internal/domain/repository"
	"duh/internal/domain/utils"
	"fmt"
	"maps"
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

func (cli *CliService) UpsertAlias(key string, value string) error {
	repo, err := cli.dbRepository.GetDefaultRepository()
	if err != nil {
		return err
	}
	repo.Aliases[key] = value
	return cli.dbRepository.UpsertRepository(*repo)
}

func (cli *CliService) RemoveAlias(key string) error {
	repo, err := cli.dbRepository.GetDefaultRepository()
	if err != nil {
		return err
	}

	delete(repo.Aliases, key)
	return cli.dbRepository.UpsertRepository(*repo)
}

func (cli *CliService) ListAliases() (map[string]string, error) {
	repos, err := cli.dbRepository.GetEnabledRepositories()
	if err != nil {
		return nil, err
	}
	entries := map[string]string{}
	for _, repo := range repos {
		maps.Copy(entries, repo.Aliases)
	}
	return entries, nil
}

func (cli *CliService) UpsertExport(key string, value string) error {
	repo, err := cli.dbRepository.GetDefaultRepository()
	if err != nil {
		return err
	}
	repo.Exports[key] = value
	return cli.dbRepository.UpsertRepository(*repo)
}

func (cli *CliService) RemoveExport(key string) error {
	repo, err := cli.dbRepository.GetDefaultRepository()
	if err != nil {
		return err
	}

	delete(repo.Exports, key)
	return cli.dbRepository.UpsertRepository(*repo)
}

func (cli *CliService) ListExports() (map[string]string, error) {
	repos, err := cli.dbRepository.GetEnabledRepositories()
	if err != nil {
		return nil, err
	}
	entries := map[string]string{}
	for _, repo := range repos {
		maps.Copy(entries, repo.Exports)
	}
	return entries, nil
}
