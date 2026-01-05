package usecase

import (
	"duh/internal/domain/entity"
	"duh/internal/domain/port"
	"duh/internal/domain/utils"
	"fmt"
	"strings"
)

type InjectUsecase struct {
	dbPort       port.DbPort
	functionPort port.FunctionPort
}

func NewInjectUsecase(dbPort port.DbPort, functionPort port.FunctionPort) *InjectUsecase {
	return &InjectUsecase{
		dbPort:       dbPort,
		functionPort: functionPort,
	}
}

func (i *InjectUsecase) GetInjectionString() (string, error) {
	enabledRepos, err := i.dbPort.GetEnabledRepositories()
	if err != nil {
		return "", err
	}
	injectionLines := []string{"alias duh_reload='eval \"$(duh inject --quiet)\"'"}
	for _, repo := range enabledRepos {
		for key, value := range repo.Aliases {
			escapedValue := utils.EscapeShellString(value)
			injectionLines = append(injectionLines, fmt.Sprintf("alias %s=\"%s\"", key, escapedValue))
		}
		for key, value := range repo.Exports {
			injectionLines = append(injectionLines, fmt.Sprintf("export %s=\"%s\"", key, value))
		}
	}
	injectionString := strings.Join(injectionLines, "\n")

	activatedScripts, _ := i.getActivatedScripts()
	for _, script := range activatedScripts {
		injectionString = fmt.Sprintf("%s\n%s", injectionString, script.DataToInject)
	}

	bonus, _ := i.dbPort.BonusInjection(enabledRepos)
	injectionString = fmt.Sprintf("%s\n%s", injectionString, bonus)
	return injectionString, nil
}

func (i *InjectUsecase) getActivatedScripts() ([]entity.Script, error) {
	scripts, err := i.functionPort.GetInternalScripts()
	if err != nil {
		return nil, err
	}
	scriptsRepos, err := i.functionPort.GetActivatedScripts()
	if err != nil {
		return nil, err
	}
	scripts = append(scripts, scriptsRepos...)
	return scripts, nil
}
