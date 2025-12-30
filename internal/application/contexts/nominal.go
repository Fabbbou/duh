package contexts

import (
	"duh/internal/application/cli"
	"duh/internal/domain/service"
	"duh/internal/infrastructure/filesystem/file_db"

	"github.com/spf13/cobra"
)

func InitCli() *cobra.Command {
	pathProvider := file_db.BasePathProvider{}
	dbRepository := file_db.NewFileDbRepository(&pathProvider, &file_db.GitConfigPathProvider{})
	cliService := service.NewCliService(dbRepository)
	return cli.BuildRootCli(cliService)
}
