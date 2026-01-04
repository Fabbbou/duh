package cli_old

import (
	"duh/internal/domain/service"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func BuildInjectSubcommand(cliService service.CliService) *cobra.Command {
	var quiet bool

	injectCmd := &cobra.Command{
		Use:   "inject",
		Short: "Inject configuration into your shell environment",
		Long:  "Generate shell commands to set up aliases and environment variables. Use --quiet for silent output suitable for sourcing.",
		Run: func(cmd *cobra.Command, args []string) {
			injection, err := cliService.Inject()
			if err != nil {
				if !quiet {
					cmd.PrintErrf("Error setting alias: %v\n", err)
				}
				return
			}

			if quiet {
				// Silent mode - output directly to stdout (for eval/sourcing)
				fmt.Fprint(os.Stdout, injection)
			} else {
				// Normal mode - show feedback to stderr, then output commands to stdout
				fmt.Fprint(os.Stdout, injection)
			}
		},
	}

	injectCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Suppress output (for silent sourcing)")

	return injectCmd
}
