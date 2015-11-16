package main

import (
	"github.com/thethingsnetwork/server-shared"
	"log"
)

var (
	consumer Consumer
	mqtt     PacketHandler
	database PacketHandler
)

func main() {
	log.Print("Jolie is ALIVE")

	err := connectConsumer()
	if err != nil {
		log.Fatalf("Failed to connect consumer: %s", err.Error())
	}

	queues, err := consumer.Consume()
	if err != nil {
		log.Fatalf("Failed to consume queues: %s", err.Error())
	}

	err = connectDatabase(queues)
	if err != nil {
		log.Fatalf("Failed to connect database: %s", err.Error())
	}

	err = connectMqtt(queues)
	if err != nil {
		log.Fatalf("Failed to connect MQTT: %s", err.Error())
	}

	select {}
}

func connectConsumer() error {
	var err error
	consumer, err = ConnectRabbitConsumer()
	if err != nil {
		return err
	}

	err = consumer.Configure()
	if err != nil {
		return err
	}

	return nil
}

func connectDatabase(queues *shared.ConsumerQueues) error {
	var err error
	database, err = ConnectMongoDatabase()
	if err != nil {
		return err
	}
	go database.Handle(queues)

	return nil
}

func connectMqtt(queues *shared.ConsumerQueues) error {
	var err error
	mqtt, err = ConnectPaho()
	if err != nil {
		return err
	}
	go mqtt.Handle(queues)

	return nil
}
