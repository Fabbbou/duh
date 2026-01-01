package repository

import "duh/internal/domain/entity"

type FunctionsRepository interface {
	// Returns the scripts from activated repositories
	// (so only when the user preferences says so)
	GetActivatedScripts() ([]entity.Script, error)
}
