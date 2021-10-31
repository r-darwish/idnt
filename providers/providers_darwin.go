//go:build darwin

package providers

func GetOsSpecificProviders() []Provider {
	return []Provider{&Brew{}}
}
