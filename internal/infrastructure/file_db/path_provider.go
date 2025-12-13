package file_db

import (
	"os"
	"path/filepath"
)

type PathProvider interface {
	GetPath() (string, error)
}

type BasePathProvider struct{}

func (bpp *BasePathProvider) GetPath() (string, error) {
	os.UserHomeDir()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".local", "share", "duh"), nil
}

type CustomPathProvider struct {
	customPath string
}

func NewCustomPathProvider(customPath string) *CustomPathProvider {
	return &CustomPathProvider{customPath: customPath}
}

func (cpp *CustomPathProvider) GetPath() (string, error) {
	return cpp.customPath, nil
}
