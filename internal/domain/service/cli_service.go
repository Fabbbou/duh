package service

import (
	"duh/internal/domain/entity"
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
	injectionLines := []string{fmt.Sprintf("alias %s=\"%s\"", "duh_reload", "eval \"$(duh inject --quiet)\"")}
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
	bonus, _ := cli.dbRepository.BonusInjection(enabledRepos)
	injectionString = fmt.Sprintf("%s\n%s", injectionString, bonus)
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

// Repository management methods
func (cli *CliService) ListRepositories() (map[string][]string, error) {
	repos, err := cli.dbRepository.GetAllRepositories()
	if err != nil {
		return nil, err
	}

	enabledRepos, err := cli.dbRepository.GetEnabledRepositories()
	if err != nil {
		return nil, err
	}

	enabledMap := make(map[string]bool)
	for _, repo := range enabledRepos {
		enabledMap[repo.Name] = true
	}

	result := map[string][]string{
		"enabled":  {},
		"disabled": {},
	}

	for _, repo := range repos {
		if enabledMap[repo.Name] {
			result["enabled"] = append(result["enabled"], repo.Name)
		} else {
			result["disabled"] = append(result["disabled"], repo.Name)
		}
	}

	return result, nil
}

func (cli *CliService) EnableRepository(repoName string) error {
	return cli.dbRepository.EnableRepository(repoName)
}

func (cli *CliService) DisableRepository(repoName string) error {
	return cli.dbRepository.DisableRepository(repoName)
}

func (cli *CliService) DeleteRepository(repoName string) error {
	return cli.dbRepository.DeleteRepository(repoName)
}

func (cli *CliService) SetDefaultRepository(repoName string) error {
	err := cli.dbRepository.EnableRepository(repoName)
	if err != nil {
		return err
	}
	return cli.dbRepository.ChangeDefaultRepository(repoName)
}

func (cli *CliService) RenameRepository(oldName, newName string) error {
	return cli.dbRepository.RenameRepository(oldName, newName)
}

func (cli *CliService) GetCurrentDefaultRepository() (string, error) {
	repo, err := cli.dbRepository.GetDefaultRepository()
	if err != nil {
		return "", err
	}
	return repo.Name, nil
}

func (cli *CliService) AddRepository(url string, name *string) error {
	repo, err := cli.dbRepository.AddRepository(url, name)
	if err != nil {
		return err
	}
	return cli.EnableRepository(repo)
}

func (cli *CliService) CreateRepository(name string) error {
	_, err := cli.dbRepository.CreateRepository(name)
	return err
}

func (cli *CliService) UpdateRepos(strategy string) (entity.RepositoryUpdateResults, error) {
	return cli.dbRepository.UpdateRepositories(strategy)
}

func (cli *CliService) EditRepo(repoName string) error {
	return cli.dbRepository.EditRepo(repoName)
}

func (cli *CliService) GetBasePath() (string, error) {
	return cli.dbRepository.GetBasePath()
}

func (cli *CliService) ListPath() ([]string, error) {
	return cli.dbRepository.ListRepoPath()
}

func (cli *CliService) PushRepository(repoName string) error {
	return cli.dbRepository.PushRepository(repoName)
}

func (cli *CliService) EditGitconfig(repoName string) error {
	return cli.dbRepository.EditGitconfig(repoName)
}
