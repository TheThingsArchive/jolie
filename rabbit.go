package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
	"time"
)

func EnsureRabbitConnection() (*amqp.Connection, error) {
	// Connects opens an AMQP connection from the credentials in the URL.
	var err error
	for i := 0; i < 20; i++ {
		conn, err := amqp.Dial(os.Getenv("AMQP_URI"))
		if err != nil {
			log.Print("Couldn't get rabbit connection")
			log.Print(err.Error())
			time.Sleep(time.Duration(2) * time.Second)
			log.Print("Retrying.....")
		} else {
			log.Print("Got rabbit connection")
			return conn, nil
		}
	}
	return nil, err
}
