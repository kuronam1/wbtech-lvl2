package main

import "fmt"

type Transport interface {
	Move() string
}

type car struct{}

func (c *car) Move() string {
	return "driving car"
}

type ship struct{}

func (s ship) Move() string {
	return "sailing ship"
}

type TransportFactory interface {
	CreateTransport() Transport
}

type CarFactory struct{}

func (c *CarFactory) CreateTransport() Transport {
	return &car{}
}

type ShipFactory struct{}

func (s *ShipFactory) CreateTransport() Transport {
	return &ship{}
}

func main() {
	var factory TransportFactory

	factory = &CarFactory{}
	transport := factory.CreateTransport()
	fmt.Println(transport.Move())

	factory = &ShipFactory{}
	transport = factory.CreateTransport()
	fmt.Println(transport.Move())
}
