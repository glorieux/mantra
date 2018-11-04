// Package hello exposes a service that returns hello
package hello

import (
	"log"

	"techmantra.io/mantra"
)

const serviceName = "hello"

// Service is an hello service
type Service struct{}

// HandleMessage handles incoming messages
func (*Service) HandleMessage(cmd mantra.Message) error {
	switch cmd.(type) {
	case Message:
		cmd.(Message) <- "Hello"
	default:
		log.Println("Unknown command")
	}
	return nil
}

func (*Service) String() string {
	return serviceName
}

// Message is a message returning "hello"
type Message chan string

// To is the message recipient
func (Message) To() string {
	return serviceName
}
