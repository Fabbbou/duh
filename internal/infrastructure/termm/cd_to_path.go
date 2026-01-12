package termm

import (
	"duh/internal/domain/errorss"
	"duh/internal/domain/utils"
	"fmt"
)

func CdTo(dirPath string) error {
	if !CommandExists("cd") {
		return &errorss.InfrastructureError{Message: fmt.Sprintf("command 'cd' does not exist, cannot go to path %s", dirPath)}
	}

	if !utils.DirectoryExists(dirPath) {
		return &errorss.InfrastructureError{Message: fmt.Sprintf("directory does not exist: %s", dirPath)}
	}

	return ExecCommand("cd", dirPath)
}
