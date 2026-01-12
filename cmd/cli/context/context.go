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
	// Initialize infrastructure
	pathProvider := &common.BasePathProvider{}
	fileHandler := &tomll.TomlFileHandler{}
	dbAdapter := file_db.NewFileDbAdapter(pathProvider, &common.GitConfigPathProvider{}, fileHandler)
	userRepository := fs_user_repository.NewFsUserRepository(fileHandler, pathProvider)
	functionRepository := fs_function_adapter.NewFSFunctionsRepository(pathProvider, userRepository)
	initDbService := file_db.NewInitDbService(pathProvider, fileHandler)

	// Initialize domain services
	aliasService := service.NewAliasService(dbAdapter)
	packageService := service.NewPackageService(dbAdapter)

	// Initialize use cases
	aliasUsecase := usecase.NewAliasUsecase(aliasService)
	exportsUsecase := usecase.NewExportsUsecase(dbAdapter)
	functionsUsecase := usecase.NewFunctionsUsecase(functionRepository)
	injectUsecase := usecase.NewInjectUsecase(dbAdapter, functionRepository)
	packageUsecase := usecase.NewPackageUsecase(packageService)
	selfUsecase := usecase.NewSelfUsecase(dbAdapter)
	initFilesystemDBUsecase := usecase.NewInitFilesystemDBUsecase(pathProvider, initDbService)

	// Initialize handlers
	initFileDBHandler := handler.NewInitFileDBHandler(initFilesystemDBUsecase)
	aliasHandler := handler.NewAliasHandler(aliasUsecase)
	exportsHandler := handler.NewExportsHandler(exportsUsecase)
	functionsHandler := handler.NewFunctionsHandler(functionsUsecase)
	injectHandler := handler.NewInjectHandler(injectUsecase)
	packageHandler := handler.NewPackageHandler(packageUsecase)
	selfHandler := handler.NewSelfHandler(selfUsecase)

	// Build and return root command
	return command.BuildRootCli(
		initFileDBHandler,
		aliasHandler,
		exportsHandler,
		functionsHandler,
		injectHandler,
		packageHandler,
		selfHandler,
	)
}
