package tomll

import (
	"duh/internal/infrastructure/filesystem/file_db"
	"strings"
)

type TomlFileHandler struct{}

// type FileHandler interface {
// 	Extension() string
// 	LoadRepositoryFile(path string) (file_db.RepositoryDto, error)
// 	SaveRepositoryFile(path string, data file_db.RepositoryDto) error
// 	LoadUserPreferenceFile(path string) (file_db.UserPreferenceDto, error)
// 	SaveUserPreferenceFile(path string, data file_db.UserPreferenceDto) error
// }

func (h *TomlFileHandler) Extension() string {
	return "toml"
}

func (h *TomlFileHandler) LoadRepositoryFile(path string) (*file_db.RepositoryDto, error) {
	repoToml, err := LoadToml[RepositoryToml](path)
	if err != nil {
		return nil, err
	}
	dto := &file_db.RepositoryDto{
		Aliases: repoToml.Aliases,
		Exports: repoToml.Exports,
		Metadata: file_db.MetadataDto{
			UrlOrigin:  repoToml.Metadata.UrlOrigin,
			NameOrigin: repoToml.Metadata.NameOrigin,
		},
	}
	return dto, nil
}

func (h *TomlFileHandler) SaveRepositoryFile(path string, data *file_db.RepositoryDto) error {
	repoToml := RepositoryToml{
		Aliases: data.Aliases,
		Exports: data.Exports,
		Metadata: MetadataMap{
			UrlOrigin:  data.Metadata.UrlOrigin,
			NameOrigin: data.Metadata.NameOrigin,
		},
	}
	return SaveToml(path, repoToml)
}

func (h *TomlFileHandler) migrateOldVersionUserPref(path string) (*UserPreferenceToml, error) {
	userPrefTomlOld, err := LoadToml[UserPreferenceTomlOld](path)
	if err != nil {
		return nil, err
	}
	newDto := &UserPreferenceToml{
		Repositories: RepositoriesPreference{
			ActivatedRepositories: []string{},
			DefaultRepositoryName: "",
		},
	}
	for key, value := range userPrefTomlOld.Repositories {
		if key == "default_repo_name" {
			newDto.Repositories.DefaultRepositoryName = value
		}
		if key == "activated_repos" {
			// Assuming the value is a comma-separated string
			var activatedRepos []string
			for _, repo := range strings.Split(value, ",") {
				activatedRepos = append(activatedRepos, repo)
			}
			newDto.Repositories.ActivatedRepositories = activatedRepos
		}
	}
	return newDto, nil
}

func (h *TomlFileHandler) LoadUserPreferenceFile(path string) (*file_db.UserPreferenceDto, error) {
	userPrefToml, err := LoadToml[UserPreferenceToml](path)
	if err != nil {
		// Try to migrate from old version
		userPrefToml, err = h.migrateOldVersionUserPref(path)
		if err != nil {
			return nil, err
		}
	}
	dto := &file_db.UserPreferenceDto{
		Repositories: file_db.RepositoriesPreferenceDto{
			ActivatedRepositories: userPrefToml.Repositories.ActivatedRepositories,
			DefaultRepositoryName: userPrefToml.Repositories.DefaultRepositoryName,
		},
	}
	return dto, nil
}

func (h *TomlFileHandler) SaveUserPreferenceFile(path string, data *file_db.UserPreferenceDto) error {
	userPrefToml := UserPreferenceToml{
		Repositories: RepositoriesPreference{
			ActivatedRepositories: data.Repositories.ActivatedRepositories,
			DefaultRepositoryName: data.Repositories.DefaultRepositoryName,
		},
	}
	return SaveToml(path, userPrefToml)
}
