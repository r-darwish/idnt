package src

import (
	"github.com/thoas/go-funk"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Pacman struct {
}

func (s Pacman) Exists() bool {
	_, err := exec.LookPath("pacman")
	return runtime.GOOS == "linux" && err == nil
}

func (s Pacman) Name() string {
	return "Pacman"
}

func (s Pacman) Packages() ([]string, error) {
	out, err := exec.Command("pacman", "-Qe").Output()
	if err != nil {
		return nil, err
	}

	var result = funk.Initial(funk.Map(strings.Split(string(out), "\n"), func(p string) string {
		return strings.Split(p, " ")[0]
	})).([]string)
	return result, nil
}

func (s Pacman) Uninstall(pkg string) error {
	cmd := exec.Command("sudo", "pacman", "-Rnsc", pkg)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
