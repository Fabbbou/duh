package handler

import (
	"duh/internal/application/usecase"
	"duh/internal/domain/entity"

	"github.com/spf13/cobra"
)

type InitFileDBHandler struct {
	initFilesystemDBUsecase *usecase.InitFilesystemDBUsecase
}

func NewInitFileDBHandler(initFilesystemDBUsecase *usecase.InitFilesystemDBUsecase) *InitFileDBHandler {
	return &InitFileDBHandler{
		initFilesystemDBUsecase: initFilesystemDBUsecase,
	}
}

func (h *InitFileDBHandler) HandleInitFileDB(cmd *cobra.Command) {
	message, err := h.initFilesystemDBUsecase.InitIfNeeded(cmd)
	if err == entity.ErrCouldNotGetPath {
		cmd.PrintErrf("Error getting Duh DB path: %v\n", err)
	} else if err == entity.ErrFSDbInitFailed {
		cmd.PrintErrf("Error checking Duh DB: %v\n", err)
	} else if message != "" {
		cmd.Print(message)
	}
}
