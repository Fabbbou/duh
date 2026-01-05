package handler

import (
	"duh/internal/application/usecase"

	"github.com/spf13/cobra"
)

type VersionHandler struct {
	versionUsecase *usecase.VersionUsecase
}

func NewVersionHandler(versionUsecase *usecase.VersionUsecase) *VersionHandler {
	return &VersionHandler{
		versionUsecase: versionUsecase,
	}
}

func (v *VersionHandler) ShowVersion(cmd *cobra.Command, args []string) {
	version := v.versionUsecase.GetVersion()
	cmd.Println(version)
}
