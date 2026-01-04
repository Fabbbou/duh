package main

import (
	"duh/cmd/cli/context"
	"fmt"
	"os"
)

func main() {
	cli := context.InitializeCLI()
	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
