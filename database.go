package main

import (
	"github.com/thethingsnetwork/server-shared"
)

type Database interface {
	//FindApplication(id string) (*Application, error)
	GetApplications() ([]*Application, error)
	//UpdateApplication(app *Application, params map[string]interface{}) error
	SaveApplication(app *Application) error

	//FindDevice(id string) (*Device, error)
	//GetDevice() ([]*Device, error)
	//UpdateDevice(device *Device, params map[string]interface{}) error
	//SaveDevice(device *Device) error

	RecordGatewayStatus(*shared.GatewayStatus) error
	RecordRxPacket(*shared.RxPacket) error
}
