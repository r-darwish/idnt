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

type BrewApp struct {
	cask bool
}

func (b *Brew) GetName() string {
	return "Brew"
}

func (b *Brew) GetApplications() ([]Application, error) {
	ctx := context.Background()
	results, err := hunch.All(ctx, func(ctx context.Context) (interface{}, error) {
		return b.getApps(true)
	},
		func(ctx context.Context) (interface{}, error) {

			return b.getApps(false)
		})
	if err != nil {
		return nil, err
	}

	var allApps []Application
	for _, result := range results {
		appGroup := result.([]Application)
		allApps = append(allApps, appGroup...)
	}

	return allApps, nil
}

func (b *Brew) getApps(cask bool) ([]Application, error) {
	var command *exec.Cmd
	if cask {
		command = exec.Command("brew", "list", "--cask")
	} else {
		command = exec.Command("brew", "leaves")
	}

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
			return b.getBrewInfo(text, cask)
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

func (b *Brew) getBrewInfo(packageName string, cask bool) (interface{}, error) {
	var command *exec.Cmd

	if cask {
		command = exec.Command("brew", "info", "--cask", packageName)
	} else {
		command = exec.Command("brew", "info", packageName)
	}

	output, err := command.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer func(output io.ReadCloser) {
		_ = output.Close()
		_ = command.Wait()
	}(output)

	err = command.Start()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(output)
	scanner.Split(bufio.ScanLines)

	extraInfo := map[string]string{}

	if cask {
		scanner.Scan()
		extraInfo["Package"] = scanner.Text()
		scanner.Scan()
		extraInfo["URL"] = scanner.Text()
		for scanner.Text() != "==> Description" {
			scanner.Scan()
		}
		scanner.Scan()
		extraInfo["Description"] = scanner.Text()
		extraInfo["Type"] = "Cask"
	} else {
		scanner.Scan()
		extraInfo["Package"] = scanner.Text()
		scanner.Scan()
		extraInfo["Description"] = scanner.Text()
		scanner.Scan()
		extraInfo["URL"] = scanner.Text()
		extraInfo["Type"] = "Formula"
	}

	return Application{
		Name:         packageName,
		Provider:     b,
		ExtendedInfo: BrewApp{cask: cask},
		ExtraInfo:    extraInfo,
	}, nil
}

func (b *Brew) RemoveApplication(application *Application) error {
	var command *exec.Cmd
	brewApp := application.ExtendedInfo.(BrewApp)

	if brewApp.cask {
		command = exec.Command("brew", "uninstall", "--cask", application.Name)
	} else {
		command = exec.Command("brew", "rmtree", "--quiet", application.Name)
	}

	pipe, err := command.StdinPipe()
	if err != nil {
		return err
	}
	defer func(pipe io.WriteCloser) {
		_ = pipe.Close()
	}(pipe)
	_, err = io.WriteString(pipe, "y\n")
	return command.Run()
}
