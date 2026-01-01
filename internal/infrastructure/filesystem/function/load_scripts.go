package function

import (
	"duh/internal/domain/entity"
	"duh/internal/domain/utils"
	"fmt"
)

// functionsDirectoryPath: path to the directory where shell scripts are stored
func GetScripts(functionsDirectoryPath string) ([]entity.Script, error) {
	if !utils.DirectoryExists(functionsDirectoryPath) {
		return nil, fmt.Errorf("could not find directory %s", functionsDirectoryPath)
	}
	scriptsPaths, err := utils.ListFilesInDirectory(functionsDirectoryPath)
	if err != nil {
		return nil, err
	}
	var scripts []entity.Script
	for _, scriptPath := range scriptsPaths {
		scriptContent, err := utils.ReadFileAsString(scriptPath)
		if err != nil {
			return nil, err
		}

		analyzer, err := GetScriptAnalysis(scriptContent)
		if err != nil {
			return nil, err
		}

		scripts = append(scripts, entity.Script{
			Name:         utils.GetFileNameWithoutExtension(scriptPath),
			PathToFile:   scriptPath,
			Functions:    analyzer.GetFunctions(),
			DataToInject: scriptContent,
			Warnings:     analyzer.GetWarnings(),
		})
	}
	return scripts, nil
}
