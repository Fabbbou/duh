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
func (s *InitDbService) Run() error {
	duhPath, err := s.pathProvider.GetPath()
	if err != nil {
		return err
	}
	// check if file ./.local/share/duh exists
	if !utils.DirectoryExists(duhPath) {
		os.MkdirAll(duhPath, os.ModePerm)
	}

	//check if ./.local/share/duh/repositories exists
	reposPath := filepath.Join(duhPath, "repositories")
	if utils.DirectoryExists(reposPath) {
		// repositories exists, no need to init
		return nil
	} else {
		os.MkdirAll(reposPath, os.ModePerm)
	}

	//check if ./.local/share/duh/repositories/local exists
	localRepoPath := filepath.Join(duhPath, "repositories", "local")
	if !utils.DirectoryExists(localRepoPath) {
		os.MkdirAll(localRepoPath, os.ModePerm)
	}

	if !utils.FileExists(filepath.Join(localRepoPath, "db.toml")) {
		file, err := os.Create(filepath.Join(localRepoPath, "db.toml"))
		if err != nil {
			return err
		}
		file.Close()
	}

	//check if ./.local/share/duh/user_preferences.toml exists
	userPrefPath := filepath.Join(duhPath, "user_preferences.toml")
	return s.InitUserPreference(userPrefPath)
}

func (svc *InitDbService) InitUserPreference(userPrefPath string) error {
	if utils.FileExists(userPrefPath) {
		return nil
	} else {
		file, err := os.Create(userPrefPath)
		if err != nil {
			return err
		}
		file.Close()
	}

	userPreferenceRepo := toml_repo.NewTomlUserPreferencesRepository(userPrefPath)
	userPrefs, err := userPreferenceRepo.Get()
	if err != nil {
		return err
	}
	if userPrefs.DefaultRepositoryName == "" {
		userPrefs.DefaultRepositoryName = "local"
	}
	if len(userPrefs.ActivatedRepositories) == 0 {
		userPrefs.ActivatedRepositories = []string{"local"}
	}
	return userPreferenceRepo.Save(userPrefs)
}
