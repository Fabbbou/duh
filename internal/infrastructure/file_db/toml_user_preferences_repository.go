package file_db

import (
	"duh/internal/domain/entity"
	"strings"
)

type TomlUserPreferencesRepository struct {
	tomlDriver *TomlDriver
}

func NewTomlUserPreferencesRepository(dbFilePath string) *TomlUserPreferencesRepository {
	return &TomlUserPreferencesRepository{
		tomlDriver: &TomlDriver{dbFilePath},
	}
}

func (r *TomlUserPreferencesRepository) Save(newVersion entity.UserPreferences) error {
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

func (r *TomlUserPreferencesRepository) Get() (entity.UserPreferences, error) {
	loaded, err := r.tomlDriver.Load()
	if err != nil {
		return entity.UserPreferences{}, err
	}
	loadedTyped, ok := loaded.(*UserPreferenceDb)
	if !ok {
		return entity.UserPreferences{}, nil
	}

	result := entity.UserPreferences{}
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
