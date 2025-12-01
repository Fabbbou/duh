package service

import (
	"duh/internal/domain/repository"
	"duh/internal/domain/utils"
	"os"
	"path/filepath"
)

type StartupService struct {
	pathProvider  PathProvider
	dbRepoFactory repository.DbRepositoryFactory
}

func NewStartupService(pathProvider PathProvider, dbRepoFactory repository.DbRepositoryFactory) *StartupService {
	return &StartupService{
		pathProvider:  pathProvider,
		dbRepoFactory: dbRepoFactory,
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
	if !utils.FileExists(userPrefPath) {
		file, err := os.Create(userPrefPath)
		if err != nil {
			return err
		}
		file.Close()
		userPrefRepo, err := s.dbRepoFactory.NewDbRepository(userPrefPath)
		if err != nil {
			return err
		}
		userPreferenceService := NewUserPreferenceService(userPrefRepo)
		userPreferenceService.InitUserPreference()
	}
	return nil
}
