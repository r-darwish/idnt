package providers

type Application struct {
	Name         string
	Provider     Provider
	ExtendedInfo interface{}
}

type Provider interface {
	GetApplications() ([]Application, error)
	RemoveApplication(application *Application) error
}
