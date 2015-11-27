package main

import (
	"github.com/thethingsnetwork/server-shared"
)

type PacketHandler interface {
	Configure() error
	HandleStatus(*shared.GatewayStatus)
	HandlePacket(*shared.RxPacket)
}
