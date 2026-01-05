package command

import (
	"duh/internal/interfaces/cli/handler"

	"github.com/spf13/cobra"
)

func BuildInjectCommand(injectHandler *handler.InjectHandler) *cobra.Command {
	injectCmd := &cobra.Command{
		Use:   "inject",
		Short: "Inject configuration into your shell environment",
		Long:  "Generate shell commands to set up aliases and environment variables. Use --quiet for silent output suitable for sourcing.",
		Run:   injectHandler.HandleInject,
	}

	injectCmd.Flags().Bool("quiet", false, "Silent output suitable for eval/sourcing")

	return injectCmd
}
