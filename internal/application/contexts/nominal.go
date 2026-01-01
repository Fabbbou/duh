package contexts

import (
	"duh/internal/application/cli"
	"duh/internal/domain/service"
	"duh/internal/infrastructure/filesystem/common"
	"duh/internal/infrastructure/filesystem/file_db"
	"duh/internal/infrastructure/filesystem/tomll"

	"github.com/spf13/cobra"
)

func InitCli() *cobra.Command {
	pathProvider := common.BasePathProvider{}
	dbRepository := file_db.NewFileDbRepository(&pathProvider, &common.GitConfigPathProvider{}, &tomll.TomlFileHandler{})
	cliService := service.NewCliService(dbRepository)
	return cli.BuildRootCli(cliService)
}
