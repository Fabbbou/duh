package repository

import "duh/internal/domain/entity"

type FunctionRepository interface {
	// Returns the scripts from activated repositories
	// (so only when the user preferences says so)
	GetActivatedScripts() ([]entity.Script, error)

	// Returns all scripts from all repositories
	GetAllScripts() ([]entity.Script, error)
}

type DummyFunctionRepository struct {
	Scripts          []entity.Script
	ActivatedScripts []entity.Script
	err              error
}

func (d *DummyFunctionRepository) GetActivatedScripts() ([]entity.Script, error) {
	return d.ActivatedScripts, d.err
}

func (d *DummyFunctionRepository) GetAllScripts() ([]entity.Script, error) {
	return d.Scripts, d.err
}
