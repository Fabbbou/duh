package termm

import (
	"duh/internal/domain/errorss"
	"duh/internal/domain/utils"
	"fmt"
	"os"
	"os/exec"
)

func CdTo(dirPath string) error {
	if !CommandExists("cd") {
		return &errorss.InfrastructureError{Message: fmt.Sprintf("command 'cd' does not exist, cannot go to path %s", dirPath)}
	}

	if !utils.DirectoryExists(dirPath) {
		return &errorss.InfrastructureError{Message: fmt.Sprintf("directory does not exist: %s", dirPath)}
	}

	// Create and run editor command
	cmd := exec.Command("cd", dirPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return &errorss.InfrastructureError{Message: fmt.Sprintf("failed to go to directory '%s'", dirPath)}
	}

	return nil
}
