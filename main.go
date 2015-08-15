package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	log.Print("Jolie is ALIVE")

	InfluxTest()

	var err error
	db, err = NewMongoSession()
	if err != nil {
		log.Fatalf("Failed to connect storage: %s", err.Error())
	}

	go http.ListenAndServe(":8080", Api())

	consumer, err := connectConsumer()
	if err != nil {
		log.Fatalf("Failed to connect consumer: %s", err.Error())
	}

	queues, err := consumer.Consume()
	for {
		select {
		case gatewayStatus := <-queues.GatewayStatuses:
			log.Printf("Gateway status: %#v", gatewayStatus)
		case rxPacket := <-queues.RxPackets:
			log.Printf("RX packet: %#v", rxPacket)
		}
	}
}

func connectConsumer() (Consumer, error) {
	consumer, err := ConnectRabbitConsumer()
	if err != nil {
		return nil, err
	}

	err = consumer.Configure()
	if err != nil {
		return nil, err
	}

	return consumer, nil
}
