package contexts

import (
	"duh/internal/application/cli"
	"duh/internal/domain/service"
	"duh/internal/infrastructure/filesystem/common"
	"duh/internal/infrastructure/filesystem/file_db"

	"github.com/spf13/cobra"
)

func InitCli() *cobra.Command {
	pathProvider := common.BasePathProvider{}
	dbRepository := file_db.NewFileDbRepository(&pathProvider, &common.GitConfigPathProvider{})
	cliService := service.NewCliService(dbRepository)
	return cli.BuildRootCli(cliService)
}
