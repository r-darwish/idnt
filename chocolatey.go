package main

import (
	"github.com/thoas/go-funk"
	"os"
	"os/exec"
	"strings"
)

type Chocolatey struct {
}

func (s Chocolatey) Name() string {
	return "Chocolatey"
}

func (s Chocolatey) Packages() ([]string, error) {
	out, err := exec.Command("choco", "list", "--localonly", "-r").Output()
	if err != nil {
		return nil, err
	}
	var result = funk.Map(strings.Split(string(out), "\n"), func(e string) string {
		return strings.Split(e, "|")[0]
	})
	return funk.Initial(result).([]string), nil
}

func (s Chocolatey) Uninstall(pkg string) error {
	cmd := exec.Command("choco", "uninstall", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
