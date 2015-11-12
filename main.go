package main

import (
	"log"
	"net/http"

	"github.com/thethingsnetwork/server-shared"
)

var (
	consumer Consumer
	database Database
	handlers []PacketHandler = make([]PacketHandler, 0)
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

	go http.ListenAndServe(":8080", Api())

	select {}
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
