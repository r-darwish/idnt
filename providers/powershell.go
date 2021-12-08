package providers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
)

type powershellModule struct {
	Name        string
	Description string
	ProjectUri  string
	Version     string
}

type Powershell struct {
}

func (p Powershell) GetApplications() ([]Application, error) {
	var result []Application

	if !powershellExists() {
		return result, nil
	}

	command := exec.Command("pwsh", "-c", "Get-InstalledModule | ConvertTo-Json -Depth 5")
	output, err := command.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer func(output io.ReadCloser) {
		_ = output.Close()
	}(output)

	err = command.Start()
	if err != nil {
		return nil, err
	}

	modulesOutput, err := ioutil.ReadAll(output)
	if err != nil {
		return nil, err
	}

	var parsedModules []powershellModule
	err = json.Unmarshal(modulesOutput, &parsedModules)
	if err != nil {
		return nil, err
	}

	for _, module := range parsedModules {
		var app Application
		app.Name = module.Name
		app.ExtraInfo = map[string]string{}
		app.ExtraInfo["Description"] = module.Description
		app.ExtraInfo["URL"] = module.ProjectUri
		app.ExtraInfo["Version"] = module.Version
		app.Provider = p
		result = append(result, app)
	}

	return result, nil

}

func powershellExists() bool {
	command := exec.Command("pwsh", "-v")
	output, err := command.StdoutPipe()
	if err != nil {
		return false
	}
	defer func(output io.ReadCloser) {
		_ = output.Close()
	}(output)

	err = command.Start()
	if err != nil {
		return false
	}

	return true
}

func (p Powershell) RemoveApplication(application *Application) error {
	_, err := exec.Command("pwsh", "-c", fmt.Sprintf("Uninstall-Module -Force -AllVersions %s", application.Name)).Output()
	return err
}

func (p Powershell) GetName() string {
	return "Powershell"
}
