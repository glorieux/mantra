// Package greet exposes a service to be used with mantra that greets people
package greet

import (
	"context"
	"fmt"

	"glorieux.io/mantra"
	"glorieux.io/mantra/example/hello"
)

const serviceName = "greet"

// Service is a greeting service
type Service struct{}

// Serve serves the greet service
func (g *Service) Serve(ctx context.Context, msgChan <-chan mantra.Message, send mantra.SendFunc) error {
	h := make(hello.Message)
	send(h)
	fmt.Printf("%s world!\n", <-h)
	return nil
}

// Stop stops the greet service
func (*Service) Stop() error {
	return nil
}

func (*Service) String() string {
	return serviceName
}
