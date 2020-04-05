package src

import (
	"github.com/thoas/go-funk"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Brew struct {
}

func (s Brew) Exists() bool {
	_, err := exec.LookPath("brew")
	return runtime.GOOS == "darwin" && err == nil
}

func (s Brew) Name() string {
	return "Homebrew"
}

func (s Brew) Packages() ([]string, error) {
	out, err := exec.Command("brew", "leaves").Output()
	if err != nil {
		return nil, err
	}

	var result = funk.Filter(strings.Split(string(out), "\n"), func(e string) bool {
		return e != ""
	}).([]string)
	return result, nil
}

func (s Brew) Uninstall(pkg string) error {
	cmd := exec.Command("brew", "rmtree", pkg)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
