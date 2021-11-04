//go:build linux

package providers

import (
	"bufio"
	"context"
	"errors"
	"github.com/aaronjan/hunch"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Apt struct {
}

func (p Apt) GetApplications() ([]Application, error) {
	var result []Application
	if _, err := os.Stat("/usr/bin/apt"); errors.Is(err, os.ErrNotExist) {
		return result, nil
	}

	command := exec.Command("apt-mark", "showmanual")
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
	var funcs []hunch.Executable

	for scanner.Scan() {
		line := scanner.Text()
		funcs = append(funcs, func(ctx context.Context) (interface{}, error) {
			return p.getApp(line)
		})
	}

	ctx := context.Background()
	results, err := hunch.All(ctx, funcs...)
	if err != nil {
		return nil, err
	}

	var apps []Application

	for _, result := range results {
		apps = append(apps, result.(Application))
	}

	return apps, nil
}

func (p Apt) getApp(appName string) (interface{}, error) {
	command := exec.Command("apt", "info", appName)
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

	app := Application{Provider: p}
	app.ExtraInfo = map[string]string{}

	for scanner.Scan() {
		line := scanner.Text()
		const packageField = "Package: "
		if strings.HasPrefix(line, packageField) {
			app.Name = line[len(packageField):]
		}

		store := func(aptField, infoField string) {
			if strings.HasPrefix(line, aptField) {
				app.ExtraInfo[infoField] = line[len(aptField):]
			}
		}

		store("Version: ", "Version")
		store("Homepage: ", "URL")
		store("Description: ", "Description")
	}

	return app, nil
}

func (p Apt) RemoveApplication(application *Application) error {
	_, err := exec.Command("sudo", "apt-get", "purge", "-y", application.Name).Output()
	return err
}

func (p Apt) GetName() string {
	return "Apt"
}
