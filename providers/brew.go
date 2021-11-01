//go:build !windows

package providers

import (
	"bufio"
	"context"
	"github.com/aaronjan/hunch"
	"io"
	"os/exec"
)

type Brew struct {
}

func (b *Brew) GetName() string {
	return "Brew"
}

func (b *Brew) GetApplications() ([]Application, error) {
	command := exec.Command("brew", "leaves")
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

	var result []Application
	var funcs []hunch.Executable

	scanner := bufio.NewScanner(output)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := scanner.Text()
		funcs = append(funcs, func(ctx context.Context) (interface{}, error) {
			return b.getBrewInfo(text)
		})
	}

	ctx := context.Background()
	all, err := hunch.All(ctx, funcs...)

	for _, app := range all {
		if app != nil {
			result = append(result, app.(Application))
		}
	}

	return result, nil
}

func (b *Brew) getBrewInfo(packageName string) (interface{}, error) {
	command := exec.Command("brew", "info", packageName)
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

	scanner := bufio.NewScanner(output)
	scanner.Split(bufio.ScanLines)

	extraInfo := map[string]string{}
	scanner.Scan()
	extraInfo["Package"] = scanner.Text()
	scanner.Scan()
	extraInfo["Description"] = scanner.Text()
	scanner.Scan()
	extraInfo["URL"] = scanner.Text()

	return Application{
		Name:         packageName,
		Provider:     b,
		ExtendedInfo: nil,
		ExtraInfo:    extraInfo,
	}, nil
}

func (b *Brew) RemoveApplication(application *Application) error {
	command := exec.Command("brew", "uninstall", application.Name)
	return command.Run()
}
