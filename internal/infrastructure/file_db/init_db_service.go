package file_db

import (
	"duh/internal/domain/utils"
	"duh/internal/infrastructure/file_db/toml_repo"
	"os"
	"path/filepath"
)

type InitDbService struct {
	pathProvider PathProvider
}

func NewInitDbService(
	pathProvider PathProvider,
) *InitDbService {
	return &InitDbService{
		pathProvider: pathProvider,
	}
}

// TODO: dont force local repo if repos exists
func (s *InitDbService) Run() (bool, error) {
	hasChanged := false
	duhPath, err := s.pathProvider.GetPath()
	if err != nil {
		return hasChanged, err
	}
	// check if file ./.local/share/duh exists
	if !utils.DirectoryExists(duhPath) {
		os.MkdirAll(duhPath, os.ModePerm)
		hasChanged = true
	}

	//check if ./.local/share/duh/repositories exists
	reposPath := filepath.Join(duhPath, "repositories")
	if utils.DirectoryExists(reposPath) {
		// repositories exists, no need to init
		//check if ./.local/share/duh/user_preferences.toml exists
		userPrefPath := filepath.Join(duhPath, "user_preferences.toml")
		return s.InitUserPreference(userPrefPath)
	}

	os.MkdirAll(reposPath, os.ModePerm)
	hasChanged = true

	//check if ./.local/share/duh/repositories/local exists
	localRepoPath := filepath.Join(duhPath, "repositories", "local")
	if !utils.DirectoryExists(localRepoPath) {
		os.MkdirAll(localRepoPath, os.ModePerm)
	}

	if !utils.FileExists(filepath.Join(localRepoPath, "db.toml")) {
		file, err := os.Create(filepath.Join(localRepoPath, "db.toml"))
		if err != nil {
			return hasChanged, err
		}
		file.Close()
	}

	//check if ./.local/share/duh/user_preferences.toml exists
	userPrefPath := filepath.Join(duhPath, "user_preferences.toml")
	return s.InitUserPreference(userPrefPath)
}

func (svc *InitDbService) InitUserPreference(userPrefPath string) (bool, error) {
	if utils.FileExists(userPrefPath) {
		return false, nil
	} else {
		file, err := os.Create(userPrefPath)
		if err != nil {
			return true, err
		}
		file.Close()
	}

	userPrefs, err := toml_repo.LoadUserPreferences(userPrefPath)
	if err != nil {
		return false, err
	}
	hasChanged := false
	if userPrefs.GetDefaultRepositoryName() == "" {
		hasChanged = true
		userPrefs.SetDefaultRepositoryName("local")
	}
	if len(userPrefs.GetActivatedRepositories()) == 0 {
		hasChanged = true
		userPrefs.SetActivatedRepositories([]string{"local"})
	}
	return hasChanged, toml_repo.SaveToml(userPrefPath, *userPrefs)
}
