package context

import (
	"duh/internal/application/usecase"
	"duh/internal/domain/service"
	"duh/internal/infrastructure/filesystem/common"
	"duh/internal/infrastructure/filesystem/file_db"
	"duh/internal/infrastructure/filesystem/fs_function_adapter"
	"duh/internal/infrastructure/filesystem/fs_user_repository"
	"duh/internal/infrastructure/filesystem/tomll"
	"duh/internal/interfaces/cli/command"
	"duh/internal/interfaces/cli/handler"

	"github.com/spf13/cobra"
)

func InitializeCLI() *cobra.Command {
	pathProvider := common.BasePathProvider{}
	fileHandler := &tomll.TomlFileHandler{}
	dbAdapter := file_db.NewFileDbAdapter(&pathProvider, &common.GitConfigPathProvider{}, fileHandler)
	userRepository := fs_user_repository.NewFsUserRepository(fileHandler, &pathProvider)
	functionRepository := fs_function_adapter.NewFSFunctionsRepository(&pathProvider, userRepository)
	cliService := service.NewCliService(dbAdapter, functionRepository)

	aliasUsecase := usecase.NewAliasUsecase(dbAdapter)
	aliasHandler := handler.NewAliasHandler(aliasUsecase)

	return command.BuildRootCli(cliService, aliasHandler)
}
