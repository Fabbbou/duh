package port

import "duh/internal/domain/entity"

type FunctionPort interface {
	// Returns the scripts from activated repositories
	// (so only when the user preferences says so)
	GetActivatedScripts() ([]entity.Script, error)

	// Returns all scripts from all repositories
	GetAllScripts() ([]entity.Script, error)

	// Returns internal scripts embedded in the binary
	GetInternalScripts() ([]entity.Script, error)

	// Creates a script by its name
	CreateScriptByName(scriptName string) (string, error)
}

type DummyFunctionRepository struct {
	Scripts          []entity.Script
	ActivatedScripts []entity.Script
	InternalScripts  []entity.Script
	err              error
}

func (d *DummyFunctionRepository) GetActivatedScripts() ([]entity.Script, error) {
	return d.ActivatedScripts, d.err
}

func (d *DummyFunctionRepository) GetAllScripts() ([]entity.Script, error) {
	return d.Scripts, d.err
}

func (d *DummyFunctionRepository) GetInternalScripts() ([]entity.Script, error) {
	return d.InternalScripts, d.err
}

func (d *DummyFunctionRepository) CreateScriptByName(scriptName string) (*entity.Script, error) {
	script := &entity.Script{
		Name: scriptName,
	}
	return script, d.err
}
