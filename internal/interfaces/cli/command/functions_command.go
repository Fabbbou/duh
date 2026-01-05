package command

import (
	"duh/internal/interfaces/cli/handler"

	"github.com/spf13/cobra"
)

func BuildFunctionsCommand(functionsHandler *handler.FunctionsHandler) *cobra.Command {
	functionsCmd := &cobra.Command{
		Use:     "functions [subcommand]",
		Aliases: []string{"function", "func", "fn", "fun"},
		Short:   "Manage shell functions injected by duh.",
	}

	listFunctionsCmd := &cobra.Command{
		Use:   "list",
		Short: "List functions",
		Args:  cobra.RangeArgs(0, 1),
		Run:   functionsHandler.ListFunctions,
	}
	listFunctionsCmd.Flags().BoolP("all", "a", false, "List all functions (not just activated ones)")
	listFunctionsCmd.Flags().Bool("core", false, "List internal core functions")

	functionDetailsCmd := &cobra.Command{
		Use:   "info [functionName]",
		Short: "Get details of a specific function",
		Args:  cobra.ExactArgs(1),
		Run:   functionsHandler.GetFunctionInfo,
	}

	addFunction := &cobra.Command{
		Use:   "add [functionName]",
		Short: "Create a new function script with the given name",
		Args:  cobra.ExactArgs(1),
		Run:   functionsHandler.CreateFunctionScript,
	}

	functionsCmd.AddCommand(listFunctionsCmd)
	functionsCmd.AddCommand(functionDetailsCmd)
	functionsCmd.AddCommand(addFunction)

	return functionsCmd
}
