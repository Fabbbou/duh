package usecase

import "duh/internal/infrastructure/termm"

func EditFile(filePath string) error {
	return termm.EditFile(filePath)
}

func CdTo(path string) error {
	return termm.CdTo(path)
}
