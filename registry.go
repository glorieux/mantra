package mantra

import (
	"pkg.glorieux.io/mantra/internal/log"
	"pkg.glorieux.io/mantra/internal/structs"
)

// Registry is a registry of services
type registry struct {
	r map[string]*service
}

func newServiceRegistry() *registry {
	registry := &registry{
		r: make(map[string]*service),
	}
	registry.AddService(registry)
	return registry
}

func (sr *registry) Stop() error {
	return nil
}

// TODO add more error checks
func (sr *registry) AddService(service Service) {
	s := newService(service)
	log.Debugf("Adding %s service", s.String())
	sr.r[s.Address.String()] = s
	sr.r[s.Address.String()].id = rootSupervisor.Add(sr.r[s.Address.String()])
}

func (sr *registry) RemoveService(service Service) {
	serviceName := structs.Name(service)
	log.Debugf("Removing %s service", serviceName)
	internalService, exists := sr.r[serviceName]
	if !exists {
		return
	}
	err := rootSupervisor.Remove(internalService.id)
	if err != nil {
		log.Error(err)
	}
	delete(sr.r, serviceName)
}
