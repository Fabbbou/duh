package cli

import (
	"fmt"
	"os/exec"
)

type CommandProcessor struct {
	verbose bool
}

func NewCommandProcessor() *CommandProcessor {
	return &CommandProcessor{
		verbose: false,
	}
}

func (cp *CommandProcessor) Exec(cmdName string, args ...string) error {
	cmd := exec.Command(cmdName, args...)
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	if cp.verbose {
		fmt.Println(string(output))
	}
	return nil
}
