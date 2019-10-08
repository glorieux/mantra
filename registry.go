package mantra

import (
	"errors"
	"fmt"

	"github.com/thejerf/suture"
	"pkg.glorieux.io/mantra/internal/log"
	"pkg.glorieux.io/mantra/internal/strings"
)

const registryServiceName = "mantra"

var (
	ErrEmptyName   = errors.New("Service name cannot be empty")
	ErrNotAlphaNum = errors.New("Service's name can only be alpha-numeric")
)

// Registry is a registry of services
type registry struct {
	supervisor *suture.Supervisor

	r map[string]*service
}

func newServiceRegistry(supervisor *suture.Supervisor) *registry {
	registry := &registry{
		supervisor: supervisor,
		r:          make(map[string]*service),
	}
	err := registry.addService(registry)
	if err != nil {
		log.Fatal(err)
	}
	return registry
}

func (sr *registry) Serve(mux ServeMux) {
	mux.Handle("add", func(e Event) {
		err := sr.addService(e.Data.(Service))
		if err != nil {
			log.Error(err)
		}
	})
	mux.Handle("remove", func(e Event) {
		sr.removeService(e.Data.(Service))
	})
}

func (sr *registry) Stop() error {
	return nil
}

func (sr *registry) addService(service Service) error {
	if service.String() == "" {
		return ErrEmptyName
	}

	if _, ok := sr.r[sr.String()]; ok && service.String() == sr.String() {
		return fmt.Errorf("Do not use %s as service name", service.String())
	}

	if !strings.HasAlphaNumeric(service.String()) {
		return ErrNotAlphaNum
	}

	log.Debugf("Adding %s service", service.String())
	sr.r[service.String()] = newService(service)
	sr.r[service.String()].id = sr.supervisor.Add(sr.r[service.String()])
	return nil
}

func (sr *registry) removeService(service Service) {
	log.Debugf("Removing %s service", service.String())
	internalService, exists := sr.r[service.String()]
	if !exists {
		return
	}
	sr.supervisor.Remove(internalService.id)
	delete(sr.r, service.String())
}

func (sr *registry) String() string {
	return registryServiceName
}
