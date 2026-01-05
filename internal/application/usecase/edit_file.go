package usecase

import "duh/internal/infrastructure/filesystem/editor"

func EditFile(filePath string) error {
	return editor.EditFile(filePath)
}
