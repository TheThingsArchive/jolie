package main

import (
	"fmt"
	"github.com/influxdb/influxdb/client"
	"log"
	"net/url"
	"os"
)

func InfluxTest() {
	log.Println("TESTING INFLUX")
	u, err := url.Parse(
		fmt.Sprintf(
			"http://%s:%s",
			os.Getenv("INFLUXDB_URL"),
			os.Getenv("INFLUXDB_PORT"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	conf := client.Config{
		URL:      *u,
		Username: os.Getenv("INFLUXDB_USERNAME"),
		Password: os.Getenv("INFLUXDB_PASSWORD"),
	}

	con, err := client.NewClient(conf)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("PINGING INFLUX")
	dur, ver, err := con.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Happy as a Hippo! %v, %s", dur, ver)
}
