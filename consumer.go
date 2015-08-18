package main

import (
	"github.com/thethingsnetwork/server-shared"
)

type Consumer interface {
	Configure() error
	Consume() (*shared.ConsumerQueues, error)
}
