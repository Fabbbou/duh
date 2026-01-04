package cli_old

import (
	"duh/internal/infrastructure/filesystem/common"
	"duh/internal/infrastructure/filesystem/file_db"
	"duh/internal/infrastructure/filesystem/tomll"

	"github.com/spf13/cobra"
)

func CheckDuhFileDBCreated(cmd *cobra.Command) {
	pathProvider := common.BasePathProvider{}
	initDbService := file_db.NewInitDbService(&pathProvider, &tomll.TomlFileHandler{})

	path, err := pathProvider.GetPath()
	if err != nil {
		cmd.PrintErrf("Error getting Duh DB path: %v\n", err)
		return
	}

	hasChanged, err := initDbService.Check()
	if err != nil {
		cmd.PrintErrf("Error checking Duh DB: %v\n", err)
		return
	}
	if hasChanged {
		cmd.Printf("Duh config initialized in path '%s'\n", path)
		return
	}
}
