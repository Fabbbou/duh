package usecase

import (
	"duh/internal/domain/entity"
	"duh/internal/domain/port"
	"fmt"
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
	return f.functionPort.GetActivatedScripts()
}

func (f *FunctionsUsecase) GetAllFunctions() ([]entity.Script, error) {
	return f.functionPort.GetAllScripts()
}

func (f *FunctionsUsecase) GetInternalFunctions() ([]entity.Script, error) {
	return f.functionPort.GetInternalScripts()
}

func (f *FunctionsUsecase) GetScriptByFunctionName(functionName string) (*entity.Script, error) {
	scripts, err := f.functionPort.GetAllScripts()
	if err != nil {
		return nil, err
	}
	scriptsInternal, err := f.functionPort.GetInternalScripts()
	if err != nil {
		return nil, err
	}
	scripts = append(scripts, scriptsInternal...)
	for _, script := range scripts {
		for _, function := range script.Functions {
			if function.Name == functionName {
				return &script, nil
			}
		}
	}
	return nil, nil
}

// Returns a path to a newly created script with the given name
func (f *FunctionsUsecase) CreateScriptByName(scriptName string) (string, error) {
	scripts, err := f.functionPort.GetAllScripts()
	if err != nil {
		return "", err
	}
	scriptsInternal, err := f.functionPort.GetInternalScripts()
	if err != nil {
		return "", err
	}
	concatNames := []string{}
	scripts = append(scripts, scriptsInternal...)
	for _, script := range scripts {
		concatNames = append(concatNames, script.Name)
		if script.Name == scriptName {
			return "", fmt.Errorf("script with name '%s' already exists", scriptName)
		}
	}
	scriptRequire, err := f.functionPort.CreateScriptByName(scriptName)
	if err != nil {
		return "", err
	}
	return scriptRequire, nil
}
