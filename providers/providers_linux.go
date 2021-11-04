//go:build linux

package providers

func GetOsSpecificProviders() []Provider {
	return []Provider{&Pacman{}, &Apt{}}
}
