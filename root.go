package root

package main

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

func main() {
    rootCmd := &cobra.Command{
        Use:   "your-cli",
        Short: "A modern Go CLI tool",
    }

    rootCmd.AddCommand(&cobra.Command{
        Use:   "greet [name]",
        Short: "Say hello",
        Args:  cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Printf("Hello, %s!\n", args[0])
        },
    })

    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
