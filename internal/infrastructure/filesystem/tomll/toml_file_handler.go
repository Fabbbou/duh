package tomll

import (
	"duh/internal/infrastructure/filesystem/file_db"
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
	return toRepositoryDto(repoToml), nil
}

func (h *TomlFileHandler) SaveRepositoryFile(path string, data *file_db.RepositoryDto) error {
	repoToml := toRepositoryToml(data)
	return SaveToml(path, repoToml)
}

func (h *TomlFileHandler) LoadUserPreferenceFile(path string) (*file_db.UserPreferenceDto, error) {
	userPrefToml, err := LoadToml[UserPreferenceToml](path)
	if err != nil {
		// Try to migrate from old version
		userPrefToml, err = migrateOldVersionUserPref(path)
		if err != nil {
			return nil, err
		}
	}
	return toUserPreferenceDto(userPrefToml), nil
}

func (h *TomlFileHandler) SaveUserPreferenceFile(path string, data *file_db.UserPreferenceDto) error {
	userPrefToml := toUserPreferenceToml(data)
	return SaveToml(path, userPrefToml)
}
