package mantra

import (
	"github.com/thejerf/suture"
	"pkg.glorieux.io/mantra/internal/log"
	"pkg.glorieux.io/mantra/internal/strings"
)

const registryServiceName = "mantra"

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
	registry.addService(registry)
	return registry
}

func (sr *registry) Serve(mux ServeMux) {
	mux.Handle("add", func(e Event) {
		sr.addService(e.Data.(Service))
	})
	mux.Handle("remove", func(e Event) {
		sr.removeService(e.Data.(Service))
	})
}

func (sr *registry) Stop() error {
	return nil
}

func (sr *registry) addService(service Service) {
	// TODO: Check Service names constrains and conflicts
	if service.String() == "" {
		log.Fatal("Service name cannot be empty")
	}

	if _, ok := sr.r[sr.String()]; ok && service.String() == sr.String() {
		log.Fatalf("Do not use %s as service name", service.String())
	}

	if !strings.HasAlphaNumeric(service.String()) {
		log.Fatalf("Service's name can only be alpha-numeric: %s", service.String())
	}

	log.Debugf("Adding %s service", service.String())
	sr.r[service.String()] = newService(service)
	sr.r[service.String()].id = sr.supervisor.Add(sr.r[service.String()])
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
