package toml_repo

import "strings"

type RepositoryToml struct {
	Aliases map[string]string `toml:"aliases"`
	Exports map[string]string `toml:"exports"`
}

// repositories fields:
// - activated_repos: comma separated list of activated repositories
// - default_repo_name: name of the default repository
type UserPreferenceToml struct {
	Repositories map[string]string `toml:"repositories"`
}

func (u *UserPreferenceToml) GetActivatedRepositories() []string {
	for key, value := range u.Repositories {
		if key == "activated_repos" {
			if len(value) > 0 {
				return strings.Split(value, ",")
			} else {
				return []string{}
			}
		}
	}
	return []string{}
}

func (u *UserPreferenceToml) GetDefaultRepositoryName() string {
	for key, value := range u.Repositories {
		if key == "default_repo_name" {
			return value
		}
	}
	return ""
}

func (u *UserPreferenceToml) SetDefaultRepositoryName(name string) {
	u.Repositories["default_repo_name"] = name
}

func (u *UserPreferenceToml) SetActivatedRepositories(repos []string) {
	u.Repositories["activated_repos"] = strings.Join(repos, ",")
}

func LoadRepository(filePath string) (*RepositoryToml, error) {
	return LoadToml[RepositoryToml](filePath)
}

func LoadUserPreferences(filePath string) (*UserPreferenceToml, error) {
	return LoadToml[UserPreferenceToml](filePath)
}
