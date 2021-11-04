//go:build linux

package providers

import (
	"bufio"
	"errors"
	"io"
	"os"
	"os/exec"
)

type Pacman struct {
}

func (p Pacman) GetApplications() ([]Application, error) {
	var result []Application
	const fieldNameLength = 18
	if _, err := os.Stat("/usr/bin/pacman"); errors.Is(err, os.ErrNotExist) {
		return result, nil
	}

	command := exec.Command("pacman", "-Qei")
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
		application := Application{
			Provider: &p,
		}
		application.Name = scanner.Text()[fieldNameLength:]
		application.ExtraInfo = map[string]string{}

		scanner.Scan()
		application.ExtraInfo["Version"] = scanner.Text()[fieldNameLength:]

		scanner.Scan()
		application.ExtraInfo["Description"] = scanner.Text()[fieldNameLength:]

		scanner.Scan()
		scanner.Scan()
		application.ExtraInfo["URL"] = scanner.Text()[fieldNameLength:]

		result = append(result, application)

		for scanner.Scan() && scanner.Text() != "" {

		}
	}

	return result, nil
}

func (p Pacman) RemoveApplication(application *Application) error {
	_, err := exec.Command("sudo", "pacman", "-Rns", "--noconfirm", application.Name).Output()
	return err
}

func (p Pacman) GetName() string {
	return "Pacman"
}
