package main

import (
	"fmt"
	"github.com/influxdb/influxdb/client"
	"log"
	"net/url"
	"os"
	"time"
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
			ps.store("gateway_status",
				map[string]string{
					"eui": status.Eui,
				},
				map[string]interface{}{
					"latitude":  status.Latitude,
					"longitude": status.Longitude,
					"altitude":  status.Altitude,
				},
				status.Time)
		case packet := <-queues.RxPackets:
			log.Printf("Storing a RX packet %#v", packet)
			ps.store("rx_packets",
				map[string]string{
					"gateway_eui": packet.GatewayEui,
					"node_eui":    packet.NodeEui,
				},
				map[string]interface{}{
					"data": packet.Data,
				},
				packet.Time)
		}
	}
}

func (ps *InfluxDatabase) store(measurement string, tags map[string]string, fields map[string]interface{}, time time.Time) error {
	point := client.Point{
		Measurement: measurement,
		Tags:        tags,
		Fields:      fields,
		Time:        time,
		Precision:   "n",
	}

	bps := client.BatchPoints{
		Points:          []client.Point{point},
		Database:        os.Getenv("INFLUXDB_DBNAME"),
		RetentionPolicy: "default",
	}
	_, err := ps.conn.Write(bps)
	if err != nil {
		log.Printf("Failed to write data: %s", err.Error())
		return err
	}

	return nil
}
