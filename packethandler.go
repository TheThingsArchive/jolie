package main

type PacketHandler interface {
	Configure() error
	Handle(*ConsumerQueues)
}
