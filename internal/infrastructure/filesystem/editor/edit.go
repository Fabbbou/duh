package editor

import (
	"fmt"
	"os"
	"os/exec"
)

func EditFile(filePath string) error {
	// Find default editor
	editorCmd := FindDefaultFileEditor()

	// Create and run editor command
	cmd := exec.Command(editorCmd, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to open editor '%s': %w", editorCmd, err)
	}

	return nil
}
