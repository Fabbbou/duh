package file_db

import (
	"duh/internal/domain/constants"
	"duh/internal/domain/utils"
	"duh/internal/infrastructure/filesystem/common"
	"os"
	"path/filepath"
)

type InitDbService struct {
	pathProvider common.PathProvider
	fileHandler  common.FileHandler
}

func NewInitDbService(
	pathProvider common.PathProvider,
	fileHandler common.FileHandler,
) *InitDbService {
	return &InitDbService{
		pathProvider: pathProvider,
		fileHandler:  fileHandler,
	}
}

// TODO: dont force local repo if repos exists
func (s *InitDbService) Check() (bool, error) {
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
	reposPath := filepath.Join(duhPath, constants.PackagesDirName)
	if utils.DirectoryExists(reposPath) {
		// repositories exists, no need to init
		//check if ./.local/share/duh/constants.DuhConfigFileName.[ext] exists
		userPrefPath := filepath.Join(duhPath, constants.DuhConfigFileName+"."+s.fileHandler.Extension())
		return s.InitUserPreference(userPrefPath)
	}

	os.MkdirAll(reposPath, os.ModePerm)
	hasChanged = true

	//check if ./.local/share/duh/repositories/local exists
	localRepoPath := filepath.Join(duhPath, constants.PackagesDirName, "local")
	if !utils.DirectoryExists(localRepoPath) {
		os.MkdirAll(localRepoPath, os.ModePerm)
	}

	if !utils.FileExists(filepath.Join(localRepoPath, constants.PackageDbFileName+"."+s.fileHandler.Extension())) {
		file, err := os.Create(filepath.Join(localRepoPath, constants.PackageDbFileName+"."+s.fileHandler.Extension()))
		if err != nil {
			return hasChanged, err
		}
		file.Close()
	}

	//check if ./.local/share/duh/constants.DuhConfigFileName.[ext] exists
	userPrefPath := filepath.Join(duhPath, constants.DuhConfigFileName+"."+s.fileHandler.Extension())
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

	userPrefs, err := svc.fileHandler.LoadUserPreferenceFile(userPrefPath)
	if err != nil {
		return false, err
	}
	hasChanged := false
	if userPrefs.Repositories.DefaultRepositoryName == "" {
		hasChanged = true
		userPrefs.Repositories.DefaultRepositoryName = "local"
	}
	if len(userPrefs.Repositories.ActivatedRepositories) == 0 {
		hasChanged = true
		userPrefs.Repositories.ActivatedRepositories = []string{"local"}
	}
	return hasChanged, svc.fileHandler.SaveUserPreferenceFile(userPrefPath, userPrefs)
}
