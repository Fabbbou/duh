package file_service

import (
	"duh/internal/domain/utils"
	"duh/internal/infrastructure/file_db"
	"os"
	"path/filepath"
)

type StartupService struct {
	pathProvider PathProvider
}

func NewStartupService(
	pathProvider PathProvider,
) *StartupService {
	return &StartupService{
		pathProvider: pathProvider,
	}
}

func (s *StartupService) Run() error {
	duhPath, err := s.pathProvider.GetPath()
	if err != nil {
		return err
	}
	// check if file ./.local/share/duh exists
	if !utils.DirectoryExists(duhPath) {
		os.MkdirAll(duhPath, os.ModePerm)
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

func (svc *StartupService) InitUserPreference(userPrefPath string) error {
	if !utils.FileExists(userPrefPath) {
		file, err := os.Create(userPrefPath)
		if err != nil {
			return err
		}
		file.Close()
	}

	userPreferenceRepo := file_db.NewTomlUserPreferencesRepository(userPrefPath)
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
