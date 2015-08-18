package main

import (
	"log"
	"net/http"
)

var (
	consumer Consumer
	database Database
	handlers []PacketHandler = make([]PacketHandler, 0)
	store    PacketStore
)

func main() {
	log.Print("Jolie is ALIVE")

	err := connectDatabase()
	if err != nil {
		log.Fatalf("Failed to connect database: %s", err.Error())
	}

	err = connectConsumer()
	if err != nil {
		log.Fatalf("Failed to connect consumer: %s", err.Error())
	}

	err = connectStore()
	if err != nil {
		log.Fatalf("Failed to connect the store: %s", err.Error())
	}

	err = addHandlers()
	if err != nil {
		log.Fatalf("Failed to add handlers: %s", err.Error())
	}

	go http.ListenAndServe(":8080", Api())

	queues, err := consumer.Consume()
	if err != nil {
		log.Fatalf("Failed to consume queues: %s", err.Error())
	}

	log.Printf("Consuming queues in %d handlers", len(handlers))
	for _, h := range handlers {
		go h.Handle(queues)
	}

	select {}
}

func connectDatabase() error {
	var err error
	database, err = ConnectMongoDatabase()
	if err != nil {
		return err
	}

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

func connectStore() error {
	influx, err := ConnectInfluxDatabase()
	if err != nil {
		return err
	}

	err = influx.Configure()
	if err != nil {
		return err
	}

	// Influx acts both as a store and as a handler
	store = influx
	handlers = append(handlers, influx)

	return nil
}

func addHandlers() error {
	// TODO: Delete after demo
	waterSensor := NewWaterSensor()
	handlers = append(handlers, waterSensor)

	return nil
}
