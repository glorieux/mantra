package mantra

import (
	"context"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/thejerf/suture"
)

const registryServiceName = "__mantra_registry__"

// Registry is a registry of services
type registry struct {
	*Address
	*Mailbox

	log        *logrus.Logger
	supervisor *suture.Supervisor

	mu        sync.Mutex
	mailboxes map[string]*Address

	r map[string]*service
}

func newServiceRegistry(supervisor *suture.Supervisor, logger *logrus.Logger) *registry {
	fmt.Println("NEW")
	registry := &registry{
		supervisor: supervisor,
		log:        logger,
		mailboxes:  make(map[string]*Address),
		r:          make(map[string]*service),
	}
	registry.addService(registry)
	registry.Mailbox = newMailbox(registryServiceName)
	registry.Address = newAddress(registryServiceName, registry.Mailbox)
	registry.mailboxes[registryServiceName] = registry.Address
	return registry
}

func (sr *registry) Serve(ctx context.Context, app Application) error {
	fmt.Println("Serve")
	sr.Receive(func(msg interface{}) {
		switch msg.(type) {
		case AddServiceMessage:
			sr.addService(msg.(AddServiceMessage))
		case RemoveServiceMessage:
			sr.removeService(msg.(RemoveServiceMessage))
		case newMailboxMessage:
			name := msg.(newMailboxMessage).name
			mailbox := newMailbox(name)
			address := newAddress(name, mailbox)
			fmt.Println("NEW MAILBOX", name)
			sr.mu.Lock()
			sr.mailboxes[name] = address
			sr.mu.Unlock()
			msg.(newMailboxMessage).mailbox <- mailbox
		case lookupMessage:
			name := msg.(lookupMessage).name
			sr.mu.Lock()
			address := sr.mailboxes[name]
			sr.mu.Unlock()
			fmt.Println("LOOKUP", name, address)
			if address == nil {
				sr.Send(msg)
				break
			}
			msg.(lookupMessage).address <- address
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

func (sr *registry) NewMailbox(name string) *Mailbox {
	mailbox := make(chan *Mailbox)
	sr.Send(newMailboxMessage{name, mailbox})
	return <-mailbox
}

func (sr *registry) Lookup(name string) *Address {
	address := make(chan *Address)
	sr.Send(lookupMessage{name, address})
	return <-address
}

func (sr *registry) addService(service AddServiceMessage) {
	fmt.Printf("Adding %s service\n", service.String())
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

type newMailboxMessage struct {
	name    string
	mailbox chan *Mailbox
}

type lookupMessage struct {
	name    string
	address chan *Address
}
