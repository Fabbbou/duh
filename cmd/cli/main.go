package main

import (
	"duh/internal/application/cli"
	"fmt"
	"os"
)

func main() {
	rootCmd := cli.BuildRootCli()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
