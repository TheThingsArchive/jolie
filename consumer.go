package main

import (
	"github.com/thethingsnetwork/server-shared"
)

type Consumer interface {
	Configure() error
	Consume() (*ConsumerQueues, error)
}

type ConsumerQueues struct {
	GatewayStatuses chan *shared.GatewayStatus
	RxPackets       chan *shared.RxPacket
}
