package file_db

import (
	"duh/internal/domain/utils/gitconfig"
	"fmt"
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

type GitConfigPathProvider struct{}

func (gcpp *GitConfigPathProvider) GetPath() (string, error) {
	gitconfigPath := gitconfig.GetGitConfigUserPath()
	if gitconfigPath == "" {
		return "", fmt.Errorf("no gitconfig file found")
	}
	return gitconfigPath, nil
}
