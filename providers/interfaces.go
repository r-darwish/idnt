package providers

type Application struct {
	Name         string
	Provider     Provider
	ExtendedInfo interface{}
	ExtraInfo    map[string]string
}

type Provider interface {
	GetApplications() ([]Application, error)
	RemoveApplication(application *Application) error
	GetName() string
}
