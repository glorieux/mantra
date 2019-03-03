package mantra

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/thejerf/suture"
	"glorieux.io/slice"
)

const registryServiceName = "registry"

// Registry is a registry of services
type registry struct {
	supervisor *suture.Supervisor
	log        *logrus.Logger

	r     map[string]*service
	hosts []string
}

func newServiceRegistry(supervisor *suture.Supervisor, logger *logrus.Logger) *registry {
	hostname, err := os.Hostname()
	if err != nil {
		logger.Fatal(err)
	}
	registry := &registry{
		supervisor: supervisor,
		log:        logger,
		r:          make(map[string]*service),
		hosts:      []string{hostname},
	}
	registry.addService(registry)
	return registry
}

func (sr *registry) Serve(ctx context.Context, msgChan <-chan Message, send SendFunc) error {
	for message := range msgChan {
		switch message.(type) {
		case AddServiceMessage:
			if message.(AddServiceMessage).Service == nil {
				return errors.New("[AddServiceMessage] Service must be passed")
			}
			sr.addService(message.(AddServiceMessage).Service)
		case RemoveServiceMessage:
			sr.removeService(message.(RemoveServiceMessage))
		case AddHostMessage:
			sr.hosts = append(sr.hosts, string(message.(AddHostMessage)))
		case RemoveHostMessage:
			sr.removeHost(string(message.(RemoveHostMessage)))
		case StopMessage:
			sr.log.Info(string(message.(StopMessage)))
			sr.supervisor.Stop()
		default:
			return errors.New("Unknown message type")
		}
	}
	return nil
}

func (sr *registry) Stop() error {
	return nil
}

func (sr *registry) addService(service Service) {
	sr.log.Debugf("Adding %s service", service.String())
	sr.r[service.String()] = newService(service, sr.log, sr.send)
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

func (sr *registry) removeHost(host string) {
	for i, h := range sr.hosts {
		if h == host {
			slice.Remove(sr.hosts, i)
			return
		}
	}
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

// AddHostMessage adds an host to the registry
type AddHostMessage string

// To returns the message recipient
func (AddHostMessage) To() string {
	return registryServiceName
}

// RemoveHostMessage removes an host from the registry
type RemoveHostMessage string

// To returns the message recipient
func (RemoveHostMessage) To() string {
	return registryServiceName
}

// StopMessage stops the registry service
type StopMessage string

// To returns the message recipient
func (StopMessage) To() string {
	return registryServiceName
}
