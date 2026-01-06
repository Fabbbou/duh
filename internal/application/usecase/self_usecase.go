package usecase

import (
	"duh/internal/domain/port"
	"duh/internal/domain/utils/version"
	"os"
	"path/filepath"
)

type SelfUsecase struct {
	dbPort port.DbPort
}

func NewSelfUsecase(dbPort port.DbPort) *SelfUsecase {
	return &SelfUsecase{
		dbPort: dbPort,
	}
}

func (p *SelfUsecase) GetBasePath() (string, error) {
	return p.dbPort.GetBasePath()
}

func (p *SelfUsecase) GetAllPaths() ([]string, error) {
	return p.dbPort.ListRepoPath()
}

func (p *SelfUsecase) RepositoriesPath() (string, error) {
	path, err := p.dbPort.GetBasePath()
	if err != nil {
		return "", err
	}
	repoPath := filepath.Join(path, "repositories")
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return "", nil
	}
	return repoPath, nil
}

func (s *SelfUsecase) GetVersion() string {
	return version.BuildInfo()
}
