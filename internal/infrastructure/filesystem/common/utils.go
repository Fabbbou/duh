package common

import (
	"path/filepath"
)

func getUserPrefPath(pathProvider PathProvider, extension string) (string, error) {
	basePath, err := pathProvider.GetPath()
	if err != nil {
		return "", err
	}
	fileName := "user_preferences." + extension
	return filepath.Join(basePath, fileName), nil
}

func GetUserPreferences(pathProvider PathProvider, fileHandler FileHandler) (*UserPreferenceDto, error) {
	userPrefPath, err := getUserPrefPath(pathProvider, fileHandler.Extension())
	if err != nil {
		return nil, err
	}
	return fileHandler.LoadUserPreferenceFile(userPrefPath)
}
