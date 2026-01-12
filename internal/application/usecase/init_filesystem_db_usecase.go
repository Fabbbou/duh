package usecase

import (
	"duh/internal/domain/errorss"
	"duh/internal/infrastructure/filesystem/common"
	"duh/internal/infrastructure/filesystem/file_db"
	"fmt"

	"github.com/spf13/cobra"
)

type InitFilesystemDBUsecase struct {
	pathProvider  common.PathProvider
	initDbService *file_db.InitDbService
}

func NewInitFilesystemDBUsecase(
	pathProvider common.PathProvider,
	initDbService *file_db.InitDbService,
) *InitFilesystemDBUsecase {
	return &InitFilesystemDBUsecase{
		pathProvider:  pathProvider,
		initDbService: initDbService,
	}
}

func (i *InitFilesystemDBUsecase) InitIfNeeded(cmd *cobra.Command) (string, error) {
	path, err := i.pathProvider.GetPath()
	if err != nil {
		return "", errorss.ErrCouldNotGetPath
	}

	hasChanged, err := i.initDbService.Check()
	if err != nil {
		return "", errorss.ErrFSDbInitFailed
	}
	if hasChanged {
		return fmt.Sprintf("Duh config initialized in path '%s'\n", path), nil
	}
	return "", nil
}
