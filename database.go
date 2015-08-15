package main

type Database interface {
	//FindApplication(id string) (*Application, error)
	GetApplications() ([]*Application, error)
	//UpdateApplication(app *Application, params map[string]interface{}) error
	SaveApplication(app *Application) error

	//FindDevice(id string) (*Device, error)
	//GetDevice() ([]*Device, error)
	//UpdateDevice(device *Device, params map[string]interface{}) error
	//SaveDevice(device *Device) error
}
