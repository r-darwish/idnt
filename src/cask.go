package src

import (
	"github.com/thoas/go-funk"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Cask struct {
}

func (s Cask) Exists() bool {
	_, err := exec.LookPath("brew")
	return runtime.GOOS == "darwin" && err == nil
}

func (s Cask) Name() string {
	return "Cask"
}

func (s Cask) Packages() ([]string, error) {
	out, err := exec.Command("brew", "cask", "list").Output()
	if err != nil {
		return nil, err
	}

	var result = funk.Filter(strings.Split(string(out), "\n"), func(e string) bool {
		return e != ""
	}).([]string)
	return result, nil
}

func (s Cask) Uninstall(pkg string) error {
	cmd := exec.Command("brew", "cask", "uninstall", pkg)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
