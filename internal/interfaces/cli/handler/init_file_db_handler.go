package handler

import (
	"duh/internal/application/usecase"
	"duh/internal/domain/errors"
	"duh/internal/interfaces/cli/std"

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
	if err == errors.ErrCouldNotGetPath {
		std.Errf("Error getting Duh DB path: %v\n", err)
	} else if err == errors.ErrFSDbInitFailed {
		std.Errf("Error checking Duh DB: %v\n", err)
	} else if message != "" {
		cmd.Print(message)
	}
}
