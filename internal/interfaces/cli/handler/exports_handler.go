package handler

import (
	"duh/internal/application/usecase"
	"duh/internal/interfaces/cli/std"
	"fmt"

	"github.com/spf13/cobra"
)

type ExportsHandler struct {
	exportsUsecase *usecase.ExportsUsecase
}

func NewExportsHandler(exportsUsecase *usecase.ExportsUsecase) *ExportsHandler {
	return &ExportsHandler{
		exportsUsecase: exportsUsecase,
	}
}

func (e *ExportsHandler) SetExport(cmd *cobra.Command, args []string) {
	exportName := args[0]
	value := args[1]
	err := e.exportsUsecase.SetExport(exportName, value)
	if err != nil {
		std.Errf("Error setting export: %v\n", err)
	} else {
		fmt.Printf("Export '%s' set for value '%s'\n", exportName, value)
	}
}

func (e *ExportsHandler) UnsetExport(cmd *cobra.Command, args []string) {
	exportName := args[0]
	err := e.exportsUsecase.UnsetExport(exportName)
	if err != nil {
		std.Errf("Error removing export: %v\n", err)
	} else {
		fmt.Printf("Export '%s' removed\n", exportName)
	}
}

func (e *ExportsHandler) ListExports(cmd *cobra.Command, args []string) {
	exports, err := e.exportsUsecase.ListExports()
	if err != nil {
		std.Errf("Error listing exports: %v\n", err)
		return
	}

	if len(exports) == 0 {
		cmd.Println("No exports found")
		return
	}

	cmd.Println("Current exports:")
	for name, value := range exports {
		fmt.Printf("  %s=%s\n", name, value)
	}
}
