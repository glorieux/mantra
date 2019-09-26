// Package hello exposes a service that returns hello
package hello

import (
	"context"
	"log"

	"pkg.glorieux.io/mantra"
)

const serviceName = "hello"

// Service is an hello service
type Service struct {
	*mantra.Mailbox
}

// Serve runs the service
func (s *Service) Serve(ctx context.Context, app mantra.Application) error {
	s.Mailbox = app.NewMailbox(serviceName)

	s.Receive(func(message interface{}) {
		switch message.(type) {
		case Message:
			message.(Message) <- "Hello"
		default:
			log.Println("Unknown command")
		}
	})
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
