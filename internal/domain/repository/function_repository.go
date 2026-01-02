package repository

import "duh/internal/domain/entity"

type FunctionRepository interface {
	// Returns the scripts from activated repositories
	// (so only when the user preferences says so)
	GetActivatedScripts() ([]entity.Script, error)
}

type DummyFunctionRepository struct{}

func NewDummyFunctionRepository() *DummyFunctionRepository {
	return &DummyFunctionRepository{}
}
func (d *DummyFunctionRepository) GetActivatedScripts() ([]entity.Script, error) {
	return []entity.Script{}, nil
}
