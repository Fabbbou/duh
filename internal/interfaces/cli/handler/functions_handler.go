package handler

import (
	"duh/internal/application/usecase"
	"duh/internal/domain/entity"

	"github.com/spf13/cobra"
)

type FunctionsHandler struct {
	functionsUsecase *usecase.FunctionsUsecase
}

func NewFunctionsHandler(functionsUsecase *usecase.FunctionsUsecase) *FunctionsHandler {
	return &FunctionsHandler{
		functionsUsecase: functionsUsecase,
	}
}

func (f *FunctionsHandler) ListFunctions(cmd *cobra.Command, args []string) {
	showAll, _ := cmd.Flags().GetBool("all")
	var scripts []entity.Script
	var err error

	if showAll {
		scripts, err = f.functionsUsecase.GetAllFunctions()
	} else {
		scripts, err = f.functionsUsecase.GetActivatedFunctions()
	}

	if err != nil {
		cmd.PrintErrf("Error while listing functions: %v\n", err)
		return
	}
	for _, script := range scripts {
		cmd.Printf("Script: %s\n", script.Name)
		cmd.Printf("Path: %s\n", script.PathToFile)
		if len(script.Warnings) > 0 {
			cmd.Printf("Warnings:\n")
			for _, warning := range script.Warnings {
				cmd.Printf("  - %s\n", warning.Details)
			}
		}
		for _, fun := range script.Functions {
			cmd.Printf("  - %s\n", fun.Name)
			if len(fun.Documentation) > 0 {
				cmd.Printf("    %s\n\n", fun.Documentation[0])
			}
		}
	}
}
