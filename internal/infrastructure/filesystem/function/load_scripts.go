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
	scriptContent = removeShabangLine(scriptContent)
	script := entity.Script{
		Name:         utils.GetFileNameWithoutExtension(scriptFilePath),
		PathToFile:   scriptFilePath,
		Functions:    analyzer.GetFunctions(),
		DataToInject: scriptContent,
		Warnings:     analyzer.GetWarnings(),
	}
	return &script, nil
}

func GetScriptFromString(scriptName string, scriptContent string, pathToFile string) (*entity.Script, error) {
	analyzer, err := GetScriptAnalysis(scriptContent)
	if err != nil {
		return nil, err
	}
	scriptContent = removeShabangLine(scriptContent)
	script := entity.Script{
		Name:         scriptName,
		PathToFile:   pathToFile,
		Functions:    analyzer.GetFunctions(),
		DataToInject: scriptContent,
		Warnings:     analyzer.GetWarnings(),
	}
	return &script, nil
}

func removeShabangLine(scriptContent string) string {
	scriptContentLines := utils.SplitStringByNewLine(scriptContent)
	if len(scriptContentLines) == 0 {
		return scriptContent
	}
	if isShebangLine(scriptContentLines[0]) {
		return utils.JoinStringsWithNewLine(scriptContentLines[1:])
	}
	return scriptContent
}

func isShebangLine(line string) bool {
	return len(line) >= 2 && line[0:2] == "#!"
}
