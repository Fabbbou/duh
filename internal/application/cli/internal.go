package cli

import (
	"github.com/spf13/cobra"
)

func printEntries(cmd *cobra.Command, entries map[string]string, errorMsg string, err error) {
	if err != nil {
		cmd.PrintErrf("%s: %v\n", errorMsg, err)
		return
	}
	for key, value := range entries {
		cmd.Printf("%s='%s'\n", key, value)
	}
}
