package mantra

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/thejerf/suture"
)

const registryServiceName = "registry"

// Registry is a registry of services
type registry struct {
	supervisor *suture.Supervisor
	log        *logrus.Logger
	r          map[string]*service
}

func newServiceRegistry(supervisor *suture.Supervisor, logger *logrus.Logger) *registry {
	registry := &registry{supervisor, logger, make(map[string]*service)}
	registry.addService(registry)
	return registry
}

func (sr *registry) HandleMessage(message Message) error {
	switch message.(type) {
	case AddServiceMessage:
		if message.(AddServiceMessage).Service == nil {
			return errors.New("[AddServiceMessage] Service must be passed")
		}
		sr.addService(message.(AddServiceMessage).Service)
	case RemoveServiceMessage:
		sr.removeService(message.(RemoveServiceMessage))
	case StopMessage:
		sr.log.Info(string(message.(StopMessage)))
		sr.supervisor.Stop()
	default:
		return errors.New("Unknown message type")
	}
	return nil
}

func (sr *registry) addService(service Service) {
	sr.log.Debugf("Adding %s service", service.String())
	sr.r[service.String()] = newService(service, sr.log)
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

// SendFunc sends a message to a given service
type SendFunc func(msg Message) error

func (sr *registry) send(message Message) error {
	sr.log.Debugf("Sending %T to %s", message, message.To())
	service, ok := sr.r[message.To()]
	if !ok {
		sr.log.Warnf("Service %s is not registered", message.To())
		return fmt.Errorf("Service %s is not registered", message.To())
	}
	service.messageChan <- message
	return nil
}

func (sr *registry) String() string {
	return registryServiceName
}

// AddServiceMessage adds a service to the registry
type AddServiceMessage struct {
	Service Service
}

// To returns the message recipient
func (AddServiceMessage) To() string {
	return registryServiceName
}

// RemoveServiceMessage removes a service from the registry
type RemoveServiceMessage string

// To returns the message recipient
func (RemoveServiceMessage) To() string {
	return registryServiceName
}

// StopMessage stops the registry service
type StopMessage string

// To returns the message recipient
func (StopMessage) To() string {
	return registryServiceName
}
