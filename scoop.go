package main

import (
	"github.com/thoas/go-funk"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Scoop struct {
}

func (s Scoop) Exists() bool {
	_, err := exec.LookPath("scoop")
	return runtime.GOOS == "windows" && err == nil
}

func (s Scoop) Name() string {
	return "Scoop"
}

func (s Scoop) Packages() ([]string, error) {
	out, err := exec.Command("scoop", "list").Output()
	if err != nil {
		return nil, err
	}

	var result = funk.Filter(funk.Map(funk.Drop(strings.Split(string(out), "\n"), 2), func(p string) string {
		return strings.Split(strings.TrimSpace(p), " ")[0]
	}), func(e string) bool {
		return e != ""
	}).([]string)
	return result, nil
}

func (s Scoop) Uninstall(pkg string) error {
	cmd := exec.Command("scoop", "uninstall", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
