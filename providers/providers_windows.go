//go:build windows

package providers

func GetOsSpecificProviders() []Provider {
	return []Provider{&AppWiz{}, &Choco{}}
}
