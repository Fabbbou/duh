package tomll

import (
	"duh/internal/infrastructure/filesystem/file_db"
	"strings"
)

// toRepositoryToml converts a file_db.RepositoryDto to RepositoryToml
func toRepositoryToml(dto *file_db.RepositoryDto) RepositoryToml {
	return RepositoryToml{
		Aliases:  dto.Aliases,
		Exports:  dto.Exports,
		Metadata: MetadataMap(dto.Metadata),
	}
}

// toRepositoryDto converts a RepositoryToml to file_db.RepositoryDto
func toRepositoryDto(toml *RepositoryToml) *file_db.RepositoryDto {
	return &file_db.RepositoryDto{
		Aliases:  toml.Aliases,
		Exports:  toml.Exports,
		Metadata: file_db.MetadataDto(toml.Metadata),
	}
}

// toUserPreferenceToml converts a file_db.UserPreferenceDto to UserPreferenceToml
func toUserPreferenceToml(dto *file_db.UserPreferenceDto) UserPreferenceToml {
	return UserPreferenceToml{
		Repositories: RepositoriesPreference(dto.Repositories),
	}
}

// toUserPreferenceDto converts a UserPreferenceToml to file_db.UserPreferenceDto
func toUserPreferenceDto(toml *UserPreferenceToml) *file_db.UserPreferenceDto {
	return &file_db.UserPreferenceDto{
		Repositories: file_db.RepositoriesPreferenceDto(toml.Repositories),
	}
}

// migrateOldVersionUserPref migrates an old version UserPreferenceToml to the new version
// the only change is ActivatedRepositories that was a string comma-separated before, now is a string slice
func migrateOldVersionUserPref(path string) (*UserPreferenceToml, error) {
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
