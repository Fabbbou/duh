package file_db

import (
	"path/filepath"

	"github.com/adrg/xdg"
)

type PathProvider interface {
	GetPath() (string, error)
}

type BasePathProvider struct{}

func (bpp *BasePathProvider) GetPath() (string, error) {
	path := filepath.Join(xdg.DataHome, "duh")
	return path, nil
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
