package termm

import (
	"duh/internal/domain/errorss"
	"fmt"
	"os"
	"os/exec"
)

func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func ExecCommand(cmdName string, args ...string) error {
	// Create and run editor command
	cmd := exec.Command(cmdName, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return &errorss.InfrastructureError{
			Message: fmt.Sprintf("command %s execution error: %s", cmdName, err.Error()),
		}
	}
	return nil
}
