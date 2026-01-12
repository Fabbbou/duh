package fs_user_repository

import (
	"duh/internal/domain/constants"
	"duh/internal/infrastructure/filesystem/common"
	"path/filepath"
)

type FsUserRepository struct {
	fileHandler  common.FileHandler
	pathProvider common.PathProvider
}

func NewFsUserRepository(fileHandler common.FileHandler, pathProvider common.PathProvider) *FsUserRepository {
	return &FsUserRepository{
		fileHandler:  fileHandler,
		pathProvider: pathProvider,
	}
}

func (u *FsUserRepository) getUserPreferencesPath() (string, error) {
	basePath, err := u.pathProvider.GetPath()
	if err != nil {
		return "", err
	}
	fileName := constants.DuhConfigFileName + "." + u.fileHandler.Extension()
	return filepath.Join(basePath, fileName), nil
}

func (u *FsUserRepository) GetUserPreference() (*common.UserPreferenceDto, error) {
	userPrefPath, err := u.getUserPreferencesPath()
	if err != nil {
		return nil, err
	}
	return u.fileHandler.LoadUserPreferenceFile(userPrefPath)
}

func (u *FsUserRepository) SaveUserPreference(data *common.UserPreferenceDto) error {
	userPrefPath, err := u.getUserPreferencesPath()
	if err != nil {
		return err
	}
	return u.fileHandler.SaveUserPreferenceFile(userPrefPath, data)
}
