package main

import (
	"duh/internal/application/cli"
	"duh/internal/application/contexts"
	"fmt"
	"os"
)

func main() {

	cliService, err := contexts.InitializeContexts(nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootCmd := cli.BuildRootCli(cliService)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
