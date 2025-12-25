package main

import (
	"duh/internal/application/contexts"
	"fmt"
	"os"
)

func main() {
	cli := contexts.InitCli()
	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
