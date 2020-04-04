package main

type PackageManager interface {
	Name() string
	Packages() ([]string, error)
	Uninstall(PackageName string) error
}
