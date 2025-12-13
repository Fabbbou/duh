package toml_repo

import (
	"duh/internal/infrastructure/file_db/file_dto"
	"strings"
)

type TomlUserPreferencesRepository struct {
	tomlDriver *TomlDriver[UserPreferenceDb]
}

// could just be a static function that returns a repository entity instead
func NewTomlUserPreferencesRepository(dbFilePath string) *TomlUserPreferencesRepository {
	return &TomlUserPreferencesRepository{
		tomlDriver: &TomlDriver[UserPreferenceDb]{dbFilePath},
	}
}

func (r *TomlUserPreferencesRepository) Save(newVersion file_dto.UserPreferences) error {
	//map to toml compatible struct
	mapRepos := make(map[string]string, 0)
	mapRepos["default_repo_name"] = newVersion.DefaultRepositoryName
	mapRepos["activated_repos"] = strings.Join(newVersion.ActivatedRepositories, ",")
	userPreferencesDb := UserPreferenceDb{
		Repositories: mapRepos,
	}
	return r.tomlDriver.Save(userPreferencesDb)
}

func (r *TomlUserPreferencesRepository) GetDbPath() string {
	return r.tomlDriver.filePath
}

func (r *TomlUserPreferencesRepository) Get() (file_dto.UserPreferences, error) {
	loaded, err := r.tomlDriver.Load()
	if err != nil {
		return file_dto.UserPreferences{}, err
	}
	loadedTyped := *loaded
	result := file_dto.UserPreferences{}
	for k, v := range loadedTyped.Repositories {
		switch k {
		case "default_repo_name":
			result.DefaultRepositoryName = v
		case "activated_repos":
			if len(v) > 0 {
				result.ActivatedRepositories = strings.Split(v, ",")
			} else {
				result.ActivatedRepositories = []string{}
			}
		}
	}
	return result, nil
}
