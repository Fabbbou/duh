package context

import (
	"duh/internal/application/usecase"
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
	// Initialize infrastructure
	pathProvider := &common.BasePathProvider{}
	fileHandler := &tomll.TomlFileHandler{}
	dbAdapter := file_db.NewFileDbAdapter(pathProvider, &common.GitConfigPathProvider{}, fileHandler)
	userRepository := fs_user_repository.NewFsUserRepository(fileHandler, pathProvider)
	functionRepository := fs_function_adapter.NewFSFunctionsRepository(pathProvider, userRepository)
	initDbService := file_db.NewInitDbService(pathProvider, fileHandler)

	// Initialize use cases
	aliasUsecase := usecase.NewAliasUsecase(dbAdapter)
	exportsUsecase := usecase.NewExportsUsecase(dbAdapter)
	functionsUsecase := usecase.NewFunctionsUsecase(functionRepository)
	injectUsecase := usecase.NewInjectUsecase(dbAdapter, functionRepository)
	repositoryUsecase := usecase.NewRepositoryUsecase(dbAdapter)
	versionUsecase := usecase.NewVersionUsecase()
	pathUsecase := usecase.NewPathUsecase(dbAdapter)
	initFilesystemDBUsecase := usecase.NewInitFilesystemDBUsecase(pathProvider, initDbService)

	// Initialize handlers
	initFileDBHandler := handler.NewInitFileDBHandler(initFilesystemDBUsecase)
	aliasHandler := handler.NewAliasHandler(aliasUsecase)
	exportsHandler := handler.NewExportsHandler(exportsUsecase)
	functionsHandler := handler.NewFunctionsHandler(functionsUsecase)
	injectHandler := handler.NewInjectHandler(injectUsecase)
	repositoryHandler := handler.NewRepositoryHandler(repositoryUsecase)
	versionHandler := handler.NewVersionHandler(versionUsecase)
	pathHandler := handler.NewPathHandler(pathUsecase)

	// Build and return root command
	return command.BuildRootCli(
		initFileDBHandler,
		aliasHandler,
		exportsHandler,
		functionsHandler,
		injectHandler,
		repositoryHandler,
		versionHandler,
		pathHandler,
	)
}
