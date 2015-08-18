package demo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/thethingsnetwork/server-shared"
	"log"
	"net/http"
)

const (
	DEST_NAME = "croft.things.nonred.nl"
	DEST_PORT = 3000
)

type WaterSensor struct {
}

func NewWaterSensor() *WaterSensor {
	return &WaterSensor{}
}

func (ws *WaterSensor) Configure() error {
	return nil
}

func (ws *WaterSensor) Handle(queues *shared.ConsumerQueues) {
	for {
		select {
		case packet := <-queues.RxPackets:
			buf, err := json.Marshal(packet)
			if err != nil {
				log.Printf("Failed to marshal JSON: %s", err.Error())
				continue
			}

			addr := fmt.Sprintf("http://%s:%d", DEST_NAME, DEST_PORT)
			res, err := http.Post(addr, "application/json", bytes.NewBuffer(buf))
			if err != nil {
				log.Printf("Failed to post JSON: %s", err.Error())
				continue
			}
			defer res.Body.Close()

			log.Printf("Posted RX packet to %s with status code %d\nData:%s", addr, res.StatusCode, string(buf))
		}
	}
}
