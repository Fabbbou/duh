package handler

import (
	"duh/internal/application/usecase"
	"duh/internal/interfaces/cli/std"
	"fmt"

	"github.com/spf13/cobra"
)

type AliasHandler struct {
	aliasUsecase *usecase.AliasUsecase
}

func NewAliasHandler(aliasUsecase *usecase.AliasUsecase) *AliasHandler {
	return &AliasHandler{
		aliasUsecase: aliasUsecase,
	}
}

func (a *AliasHandler) SetAlias(cmd *cobra.Command, args []string) {
	aliasName := args[0]
	value := args[1]
	err := a.aliasUsecase.SetAlias(aliasName, value)
	if err != nil {
		std.Errf("Error setting alias: %v\n", err)
	} else {
		fmt.Printf("Alias '%s' set for command '%s'\n", aliasName, value)
	}
}

func (a *AliasHandler) UnsetAlias(cmd *cobra.Command, args []string) {
	aliasName := args[0]
	err := a.aliasUsecase.UnsetAlias(aliasName)
	if err != nil {
		std.Errf("Error removing alias: %v\n", err)
	} else {
		fmt.Printf("Alias '%s' removed\n", aliasName)
	}
}

func (a *AliasHandler) ListAliases(cmd *cobra.Command, args []string) {
	entries, err := a.aliasUsecase.ListAliases()
	if err != nil {
		std.Errf("%s: %v\n", "Error listing aliases", err)
		return
	}
	for key, value := range entries {
		fmt.Printf("%s='%s'\n", key, value)
	}
}
