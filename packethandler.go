package main

import (
	"github.com/thethingsnetwork/server-shared"
)

type PacketHandler interface {
	Configure() error
	Handle(*shared.ConsumerQueues)
}
