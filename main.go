package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

func main() {
	var packageManagers = GetPackageManager()

	fzf := exec.Command("fzf", "-m")
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

	for _, packageManager := range packageManagers {
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

		err = Uninstall(pkg, packageManagers)
		if err != nil {
			fmt.Printf("Error uninstalling %s: %s", pkg, err)
		}
	}
}

var packageRE = regexp.MustCompile("(.*) \\((\\w+)\\)")

func Uninstall(pkg string, packageManagers []PackageManager) error {
	fmt.Printf("Uninstalling %s\n", pkg)

	var submatches = packageRE.FindStringSubmatch(pkg)
	for _, packageManager := range packageManagers {
		if packageManager.Name() == submatches[2] {
			return packageManager.Uninstall(pkg)
		}
	}

	panic("Unknown package manager")
}

func GetPackageManager() []PackageManager {
	if runtime.GOOS == "windows" {
		return []PackageManager{Scoop{}}
	}

	panic("Unsupported operating system")
}
