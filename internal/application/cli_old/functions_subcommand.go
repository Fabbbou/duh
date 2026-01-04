package cli_old

import (
	"duh/internal/domain/entity"
	"duh/internal/domain/service"

	"github.com/spf13/cobra"
)

func BuildFunctionsSubcommand(cliService service.CliService) *cobra.Command {

	functionsCmd := &cobra.Command{
		Use:     "functions [subcommand]",
		Aliases: []string{"function", "func", "fn", "fun"},
		Short:   "Manage shell functions injected by duh.",
	}

	listFunctionsCmd := &cobra.Command{
		Use:   "list",
		Short: "List functions",
		Args:  cobra.RangeArgs(0, 1),

		Run: func(cmd *cobra.Command, args []string) {
			showAll, _ := cmd.Flags().GetBool("all")
			var scripts []entity.Script
			var err error

			if showAll {
				scripts, err = cliService.GetAllFunctions()
			} else {
				scripts, err = cliService.GetActivatedFunctions()
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
		},
	}

	listFunctionsCmd.Flags().BoolP("all", "a", false, "List all functions (not just activated ones)")

	functionsCmd.AddCommand(listFunctionsCmd)

	return functionsCmd
}
