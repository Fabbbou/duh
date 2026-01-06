package handler

import (
	"duh/internal/application/usecase"
	"duh/internal/interfaces/cli/std"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type InjectHandler struct {
	injectUsecase *usecase.InjectUsecase
}

func NewInjectHandler(injectUsecase *usecase.InjectUsecase) *InjectHandler {
	return &InjectHandler{
		injectUsecase: injectUsecase,
	}
}

func (i *InjectHandler) HandleInject(cmd *cobra.Command, args []string) {
	quiet, _ := cmd.Flags().GetBool("quiet")

	injection, err := i.injectUsecase.GetInjectionString()
	if err != nil {
		if !quiet {
			std.Errf("Error setting alias: %v\n", err)
		}
		return
	}

	if quiet {
		// Silent mode - output directly to stdout (for eval/sourcing)
		fmt.Fprint(os.Stdout, injection)
	} else {
		// Normal mode - output via cmd
		fmt.Printf("%s\n", injection)
	}
}
