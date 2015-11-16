package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"github.com/thethingsnetwork/server-shared"
	"log"
	"os"
	"time"
)

const (
	RABBIT_ATTEMPTS = 20
	RABBIT_EXCHANGE = "messages"
)

type RabbitConsumer struct {
	conn            *amqp.Connection
	channel         *amqp.Channel
	gatewayStatuses <-chan amqp.Delivery
	rxPackets       <-chan amqp.Delivery
}

func ConnectRabbitConsumer() (Consumer, error) {
	var err error
	for i := 0; i < RABBIT_ATTEMPTS; i++ {
		uri := os.Getenv("AMQP_URI")
		conn, err := amqp.Dial(uri)
		if err != nil {
			log.Printf("Failed to connect to %s: %s", uri, err.Error())
			time.Sleep(time.Duration(2) * time.Second)
		} else {
			consumer := &RabbitConsumer{
				conn: conn,
			}
			log.Printf("Connected to %s", uri)
			return consumer, nil
		}
	}
	return nil, err
}

func (d *RabbitConsumer) Configure() error {
	c, err := d.conn.Channel()
	if err != nil {
		log.Printf("Failed to open channel: %s", err.Error())
		return err
	}

	err = c.ExchangeDeclare(RABBIT_EXCHANGE, "topic", true, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to declare exchange: %s", err.Error())
		return err
	}

	err = c.Qos(3, 0, false)
	if err != nil {
		log.Printf("Failed to set quality of service: %s", err.Error())
		return err
	}

	d.channel = c
	return nil
}

func (d *RabbitConsumer) Consume() (*shared.ConsumerQueues, error) {
	var err error
	d.gatewayStatuses, err = d.consumeQueue("gateway.status", "gateway.status")
	if err != nil {
		log.Printf("Failed to consume gateway statuses: %s", err.Error())
		return nil, err
	}

	d.rxPackets, err = d.consumeQueue("gateway.rx", "gateway.rx")
	if err != nil {
		log.Printf("Failed to consume RX packets: %s", err.Error())
		return nil, err
	}

	queues := &shared.ConsumerQueues{
		make(chan *shared.GatewayStatus),
		make(chan *shared.RxPacket),
	}
	go d.handleDeliveries(queues)
	return queues, nil
}

func (d *RabbitConsumer) consumeQueue(queueName, bindingKey string) (<-chan amqp.Delivery, error) {
	_, err := d.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to declare queue: %s", err.Error())
		return nil, err
	}

	err = d.channel.QueueBind(queueName, bindingKey, RABBIT_EXCHANGE, false, nil)
	if err != nil {
		log.Printf("Failed to bind queue: %s", err.Error())
		return nil, err
	}

	deliveries, err := d.channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to consume: %s", err.Error())
		return nil, err
	}

	return deliveries, nil
}

func (d *RabbitConsumer) handleDeliveries(queues *shared.ConsumerQueues) {
	for {
		select {
		case delivery := <-d.gatewayStatuses:
			var status shared.GatewayStatus
			err := json.Unmarshal(delivery.Body, &status)
			if err != nil {
				log.Printf("Failed to unmarshal gateway status: %s (%q)", err.Error(), delivery.Body)
				continue
			}
			queues.GatewayStatuses <- &status
			delivery.Ack(false)
		case delivery := <-d.rxPackets:
			var rxPacket shared.RxPacket
			err := json.Unmarshal(delivery.Body, &rxPacket)
			if err != nil {
				log.Printf("Failed to unmarshal RX packet: %s (%q)", err.Error(), delivery.Body)
				continue
			}
			queues.RxPackets <- &rxPacket
			delivery.Ack(false)
		}
	}
}
