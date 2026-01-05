package usecase

import (
	"duh/internal/domain/port"
)

type PathUsecase struct {
	dbPort port.DbPort
}

func NewPathUsecase(dbPort port.DbPort) *PathUsecase {
	return &PathUsecase{
		dbPort: dbPort,
	}
}

func (p *PathUsecase) GetBasePath() (string, error) {
	return p.dbPort.GetBasePath()
}

func (p *PathUsecase) GetAllPaths() ([]string, error) {
	return p.dbPort.ListRepoPath()
}
