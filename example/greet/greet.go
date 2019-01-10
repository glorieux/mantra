// Package greet exposes a service to be used with mantra that greets people
package greet

import (
	"context"
	"fmt"
	"log"

	"glorieux.io/mantra"
	"glorieux.io/mantra/example/hello"
)

const serviceName = "greet"

// Service is a greeting service
type Service struct{}

func (g *Service) Serve(ctx context.Context, msgChan <-chan mantra.Message, send mantra.SendFunc) error {
	for msg := range msgChan {
		switch msg.(type) {
		case Message:
			h := make(hello.Message)
			send(h)
			fmt.Printf("%s %s.\n", <-h, msg.(Message).Name)
			msg.(Message).Greeted <- true
		default:
			log.Println("Unknown command")
		}
	}
	return nil
}

func (*Service) Stop() error {
	return nil
}

func (*Service) String() string {
	return serviceName
}

// Message is a message print "Hello <Message.Name>"
type Message struct {
	Name    string
	Greeted chan bool
}

// To is the message recipient
func (Message) To() string {
	return serviceName
}
