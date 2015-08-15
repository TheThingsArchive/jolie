package main

import (
	"fmt"
	"github.com/influxdb/influxdb/client"
	"log"
	"net/url"
	"os"
)

type InfluxDatabase struct {
	conn *client.Client
}

func ConnectInfluxDatabase() (*InfluxDatabase, error) {
	u, err := url.Parse(
		fmt.Sprintf(
			"http://%s:%s",
			os.Getenv("INFLUXDB_URL"),
			os.Getenv("INFLUXDB_PORT"),
		),
	)
	if err != nil {
		log.Printf("Failed to build address: %s", err.Error())
		return nil, err
	}

	conf := client.Config{
		URL:      *u,
		Username: os.Getenv("INFLUXDB_USERNAME"),
		Password: os.Getenv("INFLUXDB_PASSWORD"),
	}

	conn, err := client.NewClient(conf)
	if err != nil {
		log.Printf("Failed to create a new InfluxDB client: %s", err.Error())
		return nil, err
	}

	return &InfluxDatabase{conn}, nil
}

func (ps *InfluxDatabase) Configure() error {
	return nil
}

func (ps *InfluxDatabase) Handle(queues *ConsumerQueues) {
	for {
		select {
		case status := <-queues.GatewayStatuses:
			log.Printf("Storing a gateway status %#v", status)
			ps.storeGatewayStatus(status)
		case packet := <-queues.RxPackets:
			log.Printf("Storing a RX packet %#v", packet)
			ps.storeRxPacket(packet)
		}
	}
}

func (ps *InfluxDatabase) storeGatewayStatus(status *shared.GatewayStatus) {
	point := client.Point{
		Measurement: "gatewayStatuses",
	}
}
