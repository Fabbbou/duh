package tomll

import (
	"duh/internal/infrastructure/filesystem/common"
	"strings"
)

// toRepositoryToml converts a common.RepositoryDto to RepositoryToml
func toRepositoryToml(dto *common.RepositoryDto) RepositoryToml {
	return RepositoryToml{
		Aliases:  dto.Aliases,
		Exports:  dto.Exports,
		Metadata: MetadataMap(dto.Metadata),
	}
}

// toRepositoryDto converts a RepositoryToml to common.RepositoryDto
func toRepositoryDto(toml *RepositoryToml) *common.RepositoryDto {
	return &common.RepositoryDto{
		Aliases:  toml.Aliases,
		Exports:  toml.Exports,
		Metadata: common.MetadataDto(toml.Metadata),
	}
}

// toUserPreferenceToml converts a common.UserPreferenceDto to UserPreferenceToml
func toUserPreferenceToml(dto *common.UserPreferenceDto) UserPreferenceToml {
	return UserPreferenceToml{
		Repositories: RepositoriesPreference(dto.Repositories),
	}
}

// toUserPreferenceDto converts a UserPreferenceToml to common.UserPreferenceDto
func toUserPreferenceDto(toml *UserPreferenceToml) *common.UserPreferenceDto {
	return &common.UserPreferenceDto{
		Repositories: common.RepositoriesPreferenceDto(toml.Repositories),
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
