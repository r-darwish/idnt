//go:build !windows

package providers

import (
	"bufio"
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
	scanner := bufio.NewScanner(output)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := scanner.Text()
		result = append(result, Application{
			Name:         text,
			Provider:     b,
			ExtendedInfo: nil,
		})
	}

	return result, nil
}

func (b *Brew) RemoveApplication(application *Application) error {
	command := exec.Command("brew", "uninstall", application.Name)
	return command.Run()
}
