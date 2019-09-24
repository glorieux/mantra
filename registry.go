package mantra

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/thejerf/suture"
)

const registryServiceName = "__mantra_registry__"

// Registry is a registry of services
type registry struct {
	log        *logrus.Logger
	supervisor *suture.Supervisor

	address *Address
	mailbox *Mailbox

	mailboxes map[string]*Address

	r map[string]*service
}

func newServiceRegistry(supervisor *suture.Supervisor, logger *logrus.Logger) *registry {
	address, mailbox := newMailbox(registryServiceName)
	registry := &registry{
		supervisor: supervisor,
		log:        logger,
		address:    address,
		mailbox:    mailbox,
		mailboxes:  make(map[string]*Address),
		r:          make(map[string]*service),
	}
	registry.addService(registry)
	return registry
}

func (sr *registry) Serve(ctx context.Context, app Application) error {
	sr.mailbox.Receive(func(msg interface{}) {
		switch msg.(type) {
		case AddServiceMessage:
			sr.addService(msg.(AddServiceMessage))
		case RemoveServiceMessage:
			sr.removeService(msg.(RemoveServiceMessage))
		case StopMessage:
			sr.log.Info("Stopping")
			sr.supervisor.Stop()
		default:
			sr.log.Errorf("Unknow message: %#v", msg)
		}
	})
	return nil
}

func (sr *registry) Stop() error {
	return nil
}

func (sr *registry) NewMailbox(name string) (*Address, *Mailbox) {
	address, mailbox := newMailbox(name)
	sr.mailboxes[name] = address
	return address, mailbox
}

func (sr *registry) Lookup(name string) *Address {
	return sr.mailboxes[name]
}

func (sr *registry) addService(service AddServiceMessage) {
	sr.log.Debugf("Adding %s service", service.String())
	sr.r[service.String()] = newService(service, sr.log, sr)
	sr.r[service.String()].id = sr.supervisor.Add(sr.r[service.String()])
}

func (sr *registry) removeService(serviceName RemoveServiceMessage) {
	sr.log.Debugf("Removing %s service", serviceName)
	service, exists := sr.r[string(serviceName)]
	if !exists {
		return
	}
	sr.supervisor.Remove(service.id)
	delete(sr.r, string(serviceName))
}

func (sr *registry) String() string {
	return registryServiceName
}

// AddServiceMessage adds a service to the registry
type AddServiceMessage Service

// RemoveServiceMessage removes a service from the registry
type RemoveServiceMessage string

// StopMessage stops the registry service
type StopMessage struct{}
