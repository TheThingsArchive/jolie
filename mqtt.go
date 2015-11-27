package main

import (
	"encoding/json"
	"fmt"
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/thethingsnetwork/server-shared"
	"log"
	"os"
)

type MqttConsumer struct {
	client *MQTT.Client
}

func ConnectPaho() (PacketHandler, error) {
	uri := os.Getenv("MQTT_BROKER")
	opts := MQTT.NewClientOptions().AddBroker(uri)
	opts.SetClientID("jolie")

	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	log.Printf("Connected to %s", uri)

	return &MqttConsumer{c}, nil
}

func (c *MqttConsumer) Configure() error {
	return nil
}

func (c *MqttConsumer) HandleStatus(status *shared.GatewayStatus) {
	buffer, err := json.Marshal(status)
	if err != nil {
		log.Printf("Failed to serialize gateway status: %s", err.Error())
	}
	topic := fmt.Sprintf("gateways/%s/status", status.Eui)
	token := c.client.Publish(topic, 0, false, buffer)
	token.Wait()
	if token.Error() != nil {
		log.Printf("Failed to publish status: %s", token.Error())
	}
	log.Printf("Published gateway status to topic %s", topic)
}

func (c *MqttConsumer) HandlePacket(packet *shared.RxPacket) {
	buffer, err := json.Marshal(packet)
	if err != nil {
		log.Printf("Failed to serialize packet: %s", err.Error())
	}
	topic := fmt.Sprintf("nodes/%s/packets", packet.NodeEui)
	token := c.client.Publish(topic, 0, false, buffer)
	token.Wait()
	if token.Error() != nil {
		log.Printf("Failed to publish packet: %s", token.Error())
	}
	log.Printf("Published packet to topic %s", topic)
}
