package cli

import (
	"duh/internal/domain/service"

	"github.com/spf13/cobra"
)

func BuildInjectSubcommand(cliService service.CliService) *cobra.Command {
	injectCmd := &cobra.Command{
		Use:   "inject",
		Short: "Inject configuration into your shell environment",
		Run: func(cmd *cobra.Command, args []string) {
			injection, err := cliService.Inject()
			if err != nil {
				cmd.PrintErrf("Error setting alias: %v\n", err)
				return
			}
			cmd.Print(injection)
		},
	}
	return injectCmd
}
