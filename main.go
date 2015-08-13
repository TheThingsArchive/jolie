package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	log.Print("Jolie is ALIVE")
	db = NewMongoSession()
	go http.ListenAndServe(":8080", Api())
	time.Sleep(3 * time.Second)
	InfluxTest()

	// Connects opens an AMQP connection from the credentials in the URL.
	conn, err := EnsureRabbitConnection()
	if err != nil {
		log.Fatalf("connection.open: %s", err)
	}
	defer conn.Close()

	c, err := conn.Channel()
	if err != nil {
		log.Fatalf("channel.open: %s", err)
	}

	// We declare our topology on both the publisher and consumer to ensure they
	// are the same.  This is part of AMQP being a programmable messaging model.
	//
	// See the Channel.Publish example for the complimentary declare.
	err = c.ExchangeDeclare("messages", "topic", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("exchange.declare: %s", err)
	}

	// Establish our queue topologies that we are responsible for
	type bind struct {
		queue string
		key   string
	}

	bindings := []bind{
		bind{"stat", "stat"},
		bind{"rxpk", "rxpk"},
	}

	for _, b := range bindings {
		_, err = c.QueueDeclare(b.queue, true, false, false, false, nil)
		if err != nil {
			log.Fatalf("queue.declare: %v", err)
		}

		err = c.QueueBind(b.queue, b.key, "messages", false, nil)
		if err != nil {
			log.Fatalf("queue.bind: %v", err)
		}
	}

	// Set our quality of service.  Since we're sharing 3 consumers on the same
	// channel, we want at least 3 messages in flight.
	err = c.Qos(3, 0, false)
	if err != nil {
		log.Fatalf("basic.qos: %v", err)
	}

	// Establish our consumers that have different responsibilities.  Our first
	// two queues do not ack the messages on the server, so require to be acked
	// on the client.

	statMessages, err := c.Consume("stat", "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("basic.consume: %v", err)
	}

	rxpkMessages, err := c.Consume("rxpk", "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("basic.consume: %v", err)
	}

	for {
		select {
		case delivery, ok := <-statMessages:
			delivery.Ack(false)
			log.Print(ok)
			log.Print("stat delivery received")
			log.Printf("%#v", delivery)
			log.Printf("BODY: %s", delivery.Body)
		case delivery, ok := <-rxpkMessages:
			delivery.Ack(false)
			log.Print(ok)
			log.Print("rxpk delivery received")
			log.Printf("%#v", delivery)
			log.Printf("BODY: %s", delivery.Body)
		}
	}

	err = c.Cancel("pager", false)
	if err != nil {
		log.Fatalf("basic.cancel: %v", err)
	}
}
