package function

import (
	"duh/internal/domain/entity"
	"duh/internal/domain/utils"
	"fmt"
)

// functionsDirectoryPath: path to the directory where shell scripts are stored
func GetScripts(functionsDirectoryPath string) ([]entity.Script, error) {
	if !utils.DirectoryExists(functionsDirectoryPath) {
		return nil, ErrDirNotFound
	}
	scriptsPaths, err := utils.ListFilesInDirectory(functionsDirectoryPath)
	if err != nil {
		return nil, err
	}
	var scripts []entity.Script
	errors := []error{}
	for _, scriptPath := range scriptsPaths {
		script, err := GetScript(scriptPath)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		scripts = append(scripts, *script)
	}
	if len(errors) > 0 {
		return scripts, fmt.Errorf("errors occurred while loading scripts: %v", errors)
	}
	return scripts, nil
}

func GetScript(scriptFilePath string) (*entity.Script, error) {
	scriptContent, err := utils.ReadFileAsString(scriptFilePath)
	if err != nil {
		return nil, err
	}

	analyzer, err := GetScriptAnalysis(scriptContent)
	if err != nil {
		return nil, err
	}

	script := entity.Script{
		Name:         utils.GetFileNameWithoutExtension(scriptFilePath),
		PathToFile:   scriptFilePath,
		Functions:    analyzer.GetFunctions(),
		DataToInject: scriptContent,
		Warnings:     analyzer.GetWarnings(),
	}
	return &script, nil
}
