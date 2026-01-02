package contexts

import (
	"duh/internal/application/cli"
	"duh/internal/domain/service"
	"duh/internal/infrastructure/filesystem/common"
	"duh/internal/infrastructure/filesystem/file_db"
	fs_functions_repository "duh/internal/infrastructure/filesystem/fs_function_repository"
	"duh/internal/infrastructure/filesystem/fs_user_repository"
	"duh/internal/infrastructure/filesystem/tomll"

	"github.com/spf13/cobra"
)

func InitCli() *cobra.Command {
	pathProvider := common.BasePathProvider{}
	fileHandler := &tomll.TomlFileHandler{}
	dbRepository := file_db.NewFileDbRepository(&pathProvider, &common.GitConfigPathProvider{}, fileHandler)
	userRepository := fs_user_repository.NewFsUserRepository(fileHandler, &pathProvider)
	functionRepository := fs_functions_repository.NewFSFunctionsRepository(&pathProvider, userRepository)
	cliService := service.NewCliService(dbRepository, functionRepository)
	return cli.BuildRootCli(cliService)
}
