package usecase

import (
	"duh/internal/domain/entity"
	"duh/internal/domain/port"
)

type FunctionsUsecase struct {
	functionPort port.FunctionPort
}

func NewFunctionsUsecase(functionPort port.FunctionPort) *FunctionsUsecase {
	return &FunctionsUsecase{
		functionPort: functionPort,
	}
}

func (f *FunctionsUsecase) GetActivatedFunctions() ([]entity.Script, error) {
	scripts, err := f.functionPort.GetInternalScripts()
	if err != nil {
		return nil, err
	}
	scriptsRepos, err := f.functionPort.GetActivatedScripts()
	if err != nil {
		return nil, err
	}
	scripts = append(scripts, scriptsRepos...)
	return scripts, nil
}

func (f *FunctionsUsecase) GetAllFunctions() ([]entity.Script, error) {
	scripts, err := f.functionPort.GetInternalScripts()
	if err != nil {
		return nil, err
	}
	scriptsRepos, err := f.functionPort.GetAllScripts()
	if err != nil {
		return nil, err
	}
	scripts = append(scripts, scriptsRepos...)
	return scripts, nil
}
