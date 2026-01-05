package usecase

import (
	"duh/internal/domain/utils/version"
)

type VersionUsecase struct {
}

func NewVersionUsecase() *VersionUsecase {
	return &VersionUsecase{}
}

func (v *VersionUsecase) GetVersion() string {
	return version.BuildInfo()
}
