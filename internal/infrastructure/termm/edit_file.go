package termm

func EditFile(filePath string) error {
	// Find default editor
	editorCmd := FindDefaultFileEditor()

	return ExecCommand(editorCmd, filePath)
}
