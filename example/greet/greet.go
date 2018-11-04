// Package greet exposes a service to be used with mantra that greets people
package greet

import (
	"fmt"
	"log"

	"techmantra.io/mantra"
	"techmantra.io/mantra/example/hello"
)

const serviceName = "greet"

// Service is a greeting service
type Service struct {
	send mantra.SendFunc
}

// New returns a new geet service
func New(send mantra.SendFunc) *Service {
	return &Service{send}
}

// HandleMessage handles incoming messages
func (g *Service) HandleMessage(cmd mantra.Message) error {
	switch cmd.(type) {
	case Message:
		h := make(hello.Message)
		err := g.send(h)
		if err != nil {
			return err
		}
		fmt.Printf("%s %s.\n", <-h, cmd.(Message).Name)
		cmd.(Message).Greeted <- true
	default:
		log.Println("Unknown command")
	}
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
