// Package hello exposes a service that returns hello
package hello

import (
	"context"
	"log"

	"glorieux.io/mantra"
)

const serviceName = "hello"

// Service is an hello service
type Service struct{}

// Serve runs the service
func (*Service) Serve(ctx context.Context, msgChan <-chan mantra.Message, send mantra.SendFunc) error {
	for msg := range msgChan {
		switch msg.(type) {
		case Message:
			msg.(Message) <- "Hello"
		default:
			log.Println("Unknown command")
		}
	}
	return nil
}

// Stop stops the service
func (*Service) Stop() error {
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
