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

	showCore, _ := cmd.Flags().GetBool("core")
	if showCore {
		f.showInternalFunctions(cmd)
		return
	}

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
		if len(script.Warnings) > 0 {
			cmd.Printf("Warnings:\n")
			for _, warning := range script.Warnings {
				cmd.Printf("  - %s\n", warning.Details)
			}
		}
		displayFunctionDetails(cmd, script)
	}
}

func (f *FunctionsHandler) GetFunctionInfo(cmd *cobra.Command, args []string) {
	functionName := args[0]
	script, err := f.functionsUsecase.GetScriptByFunctionName(functionName)
	if err != nil {
		cmd.PrintErrf("Error retrieving function details: %v\n", err)
		return
	}
	if script == nil {
		cmd.Printf("Function '%s' not found\n", functionName)
		return
	}
	for _, fun := range script.Functions {
		if fun.Name != functionName {
			continue
		}

		cmd.Printf("%s()\n", fun.Name)
		if len(fun.Documentation) > 0 {
			for _, docLine := range fun.Documentation {
				cmd.Printf("  %s\n", docLine)
			}
		}
		cmd.Printf("Script: %s\n", script.PathToFile)
		cmd.Printf("\n")
	}
}

// /!\ Do not test this it will open an editor
func (f *FunctionsHandler) CreateFunctionScript(cmd *cobra.Command, args []string) {
	scriptName := args[0]
	scriptPath, err := f.functionsUsecase.CreateScriptByName(scriptName)
	if err != nil {
		cmd.PrintErrf("Error creating function script: %v\n", err)
		return
	}
	if err := usecase.EditFile(scriptPath); err != nil {
		cmd.PrintErrf("Error editing function script: %v\n", err)
		return
	}
}

func (f *FunctionsHandler) showInternalFunctions(cmd *cobra.Command) {
	scripts, err := f.functionsUsecase.GetInternalFunctions()
	if err != nil {
		cmd.PrintErrf("Error while listing internal functions: %v\n", err)
		return
	}
	for _, script := range scripts {
		displayFunctionDetails(cmd, script)
	}
}

func displayFunctionDetails(cmd *cobra.Command, script entity.Script) {
	for _, fun := range script.Functions {
		cmd.Printf("- %s()\n", fun.Name)
		if len(fun.Documentation) > 0 {
			cmd.Printf("  %s\n", fun.Documentation[0])
		}
		cmd.Printf("\n")
	}
}
