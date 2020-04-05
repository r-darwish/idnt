package main

import "github.com/thoas/go-funk"

type PackageManager interface {
	Name() string
	Packages() ([]string, error)
	Uninstall(PackageName string) error
	Exists() bool
}

var PackageManagers = GetPackageManager()

func GetPackageManager() []PackageManager {
	return funk.Filter([]PackageManager{Scoop{}, Chocolatey{}, Pacman{}, Brew{}}, func(p PackageManager) bool {
		return p.Exists()
	}).([]PackageManager)
}
