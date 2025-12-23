package service

import (
	"duh/internal/domain/entity"
	"duh/internal/domain/repository"
	"duh/internal/domain/utils"
	"fmt"
	"maps"
	"slices"
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
	aliases := map[string]string{}
	for _, repo := range repos {
		maps.Copy(aliases, repo.Aliases)
	}
	return aliases, nil
}

func (cli *CliService) AddExport(key string, value string) error {
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
	exports := map[string]string{}
	for _, repo := range repos {
		maps.Copy(exports, repo.Exports)
	}
	return exports, nil
}

//////////
// Helpers
//////////

func findRepo(targetKey string, repos []entity.Repository, fieldToSearch func(m entity.Repository) map[string]string) *entity.Repository {
	idx := slices.IndexFunc(repos, func(repo entity.Repository) bool {
		mappping := fieldToSearch(repo)
		for key := range mappping {
			if key == targetKey {
				return true
			}
		}
		return false
	})
	if idx != -1 {
		return &repos[idx]
	}
	return nil
}
