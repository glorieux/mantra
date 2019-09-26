// Package greet exposes a service to be used with mantra that greets people
package greet

import (
	"context"
	"fmt"

	"pkg.glorieux.io/mantra"
	"pkg.glorieux.io/mantra/example/hello"
)

const serviceName = "greet"

// Service is a greeting service
type Service struct{}

// Serve serves the greet service
func (g *Service) Serve(ctx context.Context, app mantra.Application) error {
	address := app.Lookup("hello")

	h := make(hello.Message)
	address.Send(h)
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
