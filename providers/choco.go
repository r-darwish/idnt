//go:build windows

package providers

import (
	"bufio"
	"errors"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Choco struct {
}

func (c Choco) GetApplications() ([]Application, error) {
	var result []Application
	if _, err := os.Stat("C:\\ProgramData\\chocolatey\\bin\\choco.exe"); errors.Is(err, os.ErrNotExist) {
		return result, nil
	}

	command := exec.Command("choco", "list", "-localonly", "--limitoutput", "--detailed")
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
	for scanner.Scan() {
		line := scanner.Text()
		splitted := strings.Split(line, "|")

		app := Application{Name: splitted[0], Provider: c}
		app.ExtraInfo = map[string]string{
			"Version": splitted[1],
		}
		result = append(result, app)
	}

	g := new(errgroup.Group)
	for i := range result {
		application := &result[i]
		g.Go(func() error {
			return fillChocoInfo(application)
		})
	}

	_ = g.Wait()

	return result, nil
}

func fillChocoInfo(application *Application) error {
	command := exec.Command("choco", "info", application.Name)
	output, err := command.StdoutPipe()
	if err != nil {
		return err
	}
	defer func(output io.ReadCloser) {
		_ = output.Close()
	}(output)

	err = command.Start()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(output)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		store := func(chocoField, idntField string) {
			if strings.HasPrefix(line, chocoField) {
				application.ExtraInfo[idntField] = line[len(chocoField):]
			}
		}

		store(" Software Site: ", "URL")
		store(" Summary: ", "Description")
	}

	return nil
}

func (c Choco) RemoveApplication(application *Application) error {
	_, err := exec.Command("choco", "uninstall", application.Name).Output()
	return err
}

func (c Choco) GetName() string {
	return "Chocolatey"
}
