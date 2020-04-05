package main

import (
	"fmt"
	"github.com/r-darwish/idnt/src"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	fzf := exec.Command("fzf", "-m", "--layout=reverse-list", "--prompt=Select applications to uninstall:")
	fzf.Stderr = os.Stderr
	stdin, err := fzf.StdinPipe()
	if err != nil {
		panic(err)
	}

	stdout, err := fzf.StdoutPipe()
	if err != nil {
		panic(err)
	}

	err = fzf.Start()
	if err != nil {
		panic(err)
	}

	for _, packageManager := range src.PackageManagers {
		packages, err := packageManager.Packages()
		if err != nil {
			fmt.Printf("Error in %s: %s\n", packageManager, err)
		}

		for _, pkg := range packages {
			_, err := fmt.Fprintf(stdin, "%s (%s)\n", pkg, packageManager.Name())
			if err != nil {
				panic(err)
			}
		}
	}

	err = stdin.Close()
	if err != nil {
		panic(err)
	}

	out, err := ioutil.ReadAll(stdout)
	if err != nil {
		panic(err)
	}

	err = fzf.Wait()
	if err != nil {
		panic(err)
	}

	for _, pkg := range strings.Split(string(out), "\n") {
		if pkg == "" {
			continue
		}

		err = Uninstall(pkg)
		if err != nil {
			fmt.Printf("Error uninstalling %s: %s", pkg, err)
		}
	}
}

var packageRE = regexp.MustCompile("(.*) \\((\\w+)\\)")

func Uninstall(pkg string) error {
	fmt.Printf("Uninstalling %s\n", pkg)

	var submatches = packageRE.FindStringSubmatch(pkg)
	for _, packageManager := range src.PackageManagers {
		if packageManager.Name() == submatches[2] {
			return packageManager.Uninstall(submatches[1])
		}
	}

	panic("Unknown package manager")
}
